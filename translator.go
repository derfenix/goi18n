package i18n

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"

	"github.com/derfenix/goi18n/internal"
)

var defaultLanguage = language.Russian

var (
	builder *catalog.Builder

	supportedLanguages    []language.Tag
	supportedLanguagesMap = map[string]struct{}{}
)

func NewExternalLoader(baseURL string, header http.Header) *internal.ExternalLoader {
	return internal.NewExternalLoader(baseURL, header)
}

func NewExternalLoaderWithClient(baseURL string, header http.Header, client *http.Client) *internal.ExternalLoader {
	return internal.NewExternalLoaderWithClient(baseURL, header, client)
}

type Translatable interface {
	Translate(ctx context.Context) string
}

func TryTranslate(ctx context.Context, obj interface{}) (string, bool) {
	if translatable, ok := obj.(Translatable); ok {
		return translatable.Translate(ctx), true
	}

	return "", false
}

func Init(fs fs.ReadDirFS) error {
	if builder != nil {
		return nil
	}

	initCatalog, err := internal.InitBuilder(fs)
	if err != nil {
		return errors.Wrap(err, "init catalog")
	}

	builder = initCatalog

	// Fill the local cache
	GetLanguages()

	return nil
}

func GetPrinter(lang language.Tag) *message.Printer {
	_ = GetLanguages()

	base, _ := lang.Base()
	if _, ok := supportedLanguagesMap[lang.String()]; !ok {
		if _, ok = supportedLanguagesMap[base.String()]; !ok {
			lang = defaultLanguage
		}
	}

	p := message.NewPrinter(lang, message.Catalog(builder))

	return p
}

func GetLanguages() []language.Tag {
	if supportedLanguages == nil {
		supportedLanguages = builder.Languages()

		for idx := range supportedLanguages {
			supportedLanguagesMap[supportedLanguages[idx].String()] = struct{}{}

			base, _ := supportedLanguages[idx].Base()
			supportedLanguagesMap[base.String()] = struct{}{}
		}
	}

	return supportedLanguages
}

func RefreshTranslations() error {
	if builder == nil {
		return nil
	}

	if err := internal.RefreshTranslations(builder); err != nil {
		return errors.WithMessage(err, "refresh translations")
	}

	return nil
}

func SetExternalBuilder(b func(builder *catalog.Builder) error) {
	internal.SetExtendBuilder(b)
}

func SetExternalLoader(loader internal.Loader) {
	internal.SetExternalLoader(loader)
}
