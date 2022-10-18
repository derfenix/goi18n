package i18n

import (
	"context"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type i18nCtxType uint8

const (
	langCtxKey i18nCtxType = iota
	printerCtxKey
)

func ContextWithLang(ctx context.Context, lang language.Tag) context.Context {
	ctx = context.WithValue(ctx, langCtxKey, lang)
	ctx = context.WithValue(ctx, printerCtxKey, GetPrinter(lang))

	return ctx
}

func PrinterFromContext(ctx context.Context) *message.Printer {
	if p, ok := ctx.Value(printerCtxKey).(*message.Printer); ok {
		return p
	}

	return GetPrinter(LanguageFromContext(ctx))
}

func Sprintf(ctx context.Context, val string, args ...interface{}) string {
	return PrinterFromContext(ctx).Sprintf(val, args...)
}

func LanguageFromContext(ctx context.Context) language.Tag {
	if lang, ok := ctx.Value(langCtxKey).(language.Tag); ok {
		return lang
	}

	return defaultLanguage
}
