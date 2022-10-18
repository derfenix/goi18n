package internal

import (
	"encoding/json"
	"io"
	"io/fs"
	"path"

	"github.com/pkg/errors"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

type Loader interface {
	Load(cat *catalog.Builder) error
}

var (
	extendBuilder func(builder *catalog.Builder) error
	extLoader     Loader
)

func SetExtendBuilder(b func(builder *catalog.Builder) error) {
	extendBuilder = b
}

func SetExternalLoader(loader Loader) {
	extLoader = loader
}

func RefreshTranslations(builder *catalog.Builder) error {
	if extLoader == nil {
		return nil
	}

	if err := extLoader.Load(builder); err != nil {
		return errors.WithMessage(err, "load translations from external")
	}

	return nil
}

type Translation struct {
	Key         string   `json:"key"`
	Description string   `json:"description"`
	Translation string   `json:"translation"`
	Plural      *plurals `json:"plural"`
}

type pluralsBase struct {
	Zero  string `json:"zero"`
	One   string `json:"one"`
	Two   string `json:"two"`
	Few   string `json:"few"`
	Many  string `json:"many"`
	Other string `json:"other"`
}

type plurals struct {
	pluralsBase
	Custom map[string]string `json:"-"`
}

func (p *plurals) UnmarshalJSON(data []byte) error {
	var (
		base pluralsBase
		val  map[string]string
	)

	if err := json.Unmarshal(data, &base); err != nil {
		return errors.Wrap(err, "unmarshal base")
	}

	if err := json.Unmarshal(data, &val); err != nil {
		return errors.Wrap(err, "unmarshal custom")
	}

	// FIXME Looks hacky, there should be a better way
	delete(val, "one")
	delete(val, "zero")
	delete(val, "two")
	delete(val, "few")
	delete(val, "may")
	delete(val, "other")

	p.pluralsBase = base
	p.Custom = val

	return nil
}

func (p *plurals) cases() (cases []interface{}) {
	capacity := 12
	if len(p.Custom) > 0 {
		capacity += len(p.Custom) * 2
	}

	cases = make([]interface{}, 0, capacity)

	if p.Custom != nil {
		for cond, val := range p.Custom {
			cases = append(cases, cond, val)
		}
	}

	if p.Zero != "" {
		cases = append(cases, plural.Zero, p.Zero)
	}

	if p.One != "" {
		cases = append(cases, plural.One, p.One)
	}

	if p.Two != "" {
		cases = append(cases, plural.Two, p.Two)
	}

	if p.Few != "" {
		cases = append(cases, plural.Few, p.Few)
	}

	if p.Many != "" {
		cases = append(cases, plural.Many, p.Many)
	}

	if p.Other != "" {
		cases = append(cases, plural.Other, p.Other)
	}

	return cases
}

func InitBuilder(fs fs.ReadDirFS) (*catalog.Builder, error) {
	cat := catalog.NewBuilder()

	if err := loadTranslations(fs, cat); err != nil {
		return nil, errors.Wrap(err, "load translations")
	}

	if extendBuilder != nil {
		if err := extendBuilder(cat); err != nil {
			return nil, errors.Wrap(err, "extend builder")
		}
	}

	return cat, nil
}

func loadTranslations(files fs.ReadDirFS, cat *catalog.Builder) error {
	dir, err := files.ReadDir("locales")
	if err != nil {
		return errors.Wrap(err, "read locales dir")
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		lang, err := language.Parse(entry.Name())
		if err != nil {
			return errors.Wrapf(err, "parse language %s", entry.Name())
		}

		subDirPath := path.Join("locales", entry.Name())

		filePath := path.Join(subDirPath, "active.json")

		reader, err := files.Open(filePath)
		if err != nil {
			return errors.Wrapf(err, "open file %s", filePath)
		}

		if err := load(reader, lang, cat); err != nil {
			if cErr := reader.Close(); cErr != nil {
				_ = cErr
			}

			return errors.Wrapf(err, "load translations from %s", filePath)
		}

		if err := reader.Close(); err != nil {
			return errors.Wrapf(err, "close file %s", filePath)
		}
	}

	if extLoader != nil {
		if err := extLoader.Load(cat); err != nil {
			return errors.WithMessage(err, "load translations from external")
		}
	}

	return nil
}

func load(r io.Reader, lang language.Tag, cat *catalog.Builder) error {
	var translations []Translation
	if err := json.NewDecoder(r).Decode(&translations); err != nil {
		return errors.Wrap(err, "decode translation")
	}

	for idx := range translations {
		trans := &translations[idx]

		switch {
		case trans.Plural != nil:
			count, format := getPlaceholders(trans.Plural.Other)

			msg := plural.Selectf(count, format, trans.Plural.cases()...)

			if err := cat.Set(lang, trans.Key, msg); err != nil {
				return errors.Wrapf(err, "set message for %s", trans.Key)
			}

		case trans.Translation != "":
			if err := cat.Set(lang, trans.Key, catalog.String(trans.Translation)); err != nil {
				return errors.Wrapf(err, "set string for %s", trans.Key)
			}
		}
	}

	return nil
}

func getPlaceholders(s string) (count int, format string) {
	for idx := 0; idx < len(s); idx++ {
		if s[idx] == '%' {
			count++

			if format == "" {
				format = s[idx : idx+1]
			}
			idx++
		}
	}

	return count, format
}
