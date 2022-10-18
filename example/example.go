package example

import (
	"context"
	"embed"
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	i18n "github.com/derfenix/goi18n"
)

//go:embed locales/**/*.json
var locales embed.FS

func Basic() {
	if err := i18n.Init(locales); err != nil {
		panic(err)
	}

	printer := i18n.GetPrinter(language.Russian)
	fmt.Println(printer.Sprintf("test", "теста"))

	printer = i18n.GetPrinter(language.English)
	fmt.Println(printer.Sprintf("test", "beer"))

	ctx := i18n.ContextWithLang(context.Background(), language.Russian)
	handler(ctx)

	ctx = i18n.ContextWithLang(context.Background(), language.English)
	handler(ctx)
}

func handler(ctx context.Context) {
	translated := i18n.Sprintf(ctx, "test plural", 2)
	fmt.Println(translated)

	lang := i18n.LanguageFromContext(ctx)
	script, _ := lang.Script()
	fmt.Println("Использован", display.Languages(language.Russian).Name(lang), "язык,", display.Scripts(language.Russian).Name(script))
}
