//go:build !i18n_extra

package i18n_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	. "github.com/derfenix/goi18n"
	"github.com/derfenix/goi18n/internal"
)

var supportedLanguages = []language.Tag{language.Russian, language.English}

func TestPrinter(t *testing.T) {
	t.Parallel()

	require.NoError(t, Init(internal.TestFS))

	t.Run("translation engaged", func(t *testing.T) {
		t.Parallel()

		t.Run("russian", func(t *testing.T) {
			t.Parallel()

			printer := GetPrinter(language.Russian)
			require.Equal(t, "Тест пива", printer.Sprintf("test", "пива"))
			assert.Equal(t, "111,223", printer.Sprint(111.223))
		})

		t.Run("russian (extended)", func(t *testing.T) {
			t.Parallel()

			printer := GetPrinter(language.MustParse("ru_RU"))
			require.Equal(t, "Тест пива", printer.Sprintf("test", "пива"))
			assert.Equal(t, "111,223", printer.Sprint(111.223))
		})

		t.Run("russian (extended kz)", func(t *testing.T) {
			t.Parallel()

			printer := GetPrinter(language.MustParse("ru_KZ"))
			require.Equal(t, "Тест пива", printer.Sprintf("test", "пива"))
			assert.Equal(t, "111,223", printer.Sprint(111.223))
		})

		t.Run("english (exactly)", func(t *testing.T) {
			t.Parallel()

			printer := GetPrinter(language.English)
			require.Equal(t, "Test of the beer", printer.Sprintf("test", "beer"))
			assert.Equal(t, "111.223", printer.Sprint(111.223))
		})

		t.Run("english (dialect)", func(t *testing.T) {
			t.Parallel()

			printer := GetPrinter(language.AmericanEnglish)
			require.Equal(t, "Test of the gun", printer.Sprintf("test", "gun"))
			assert.Equal(t, "111.223", printer.Sprint(111.223))
		})

		t.Run("unsupported language", func(t *testing.T) {
			t.Parallel()

			printer := GetPrinter(language.Albanian)
			require.Equal(t, "Тест пива", printer.Sprintf("test", "пива"))
			assert.Equal(t, "111,223", printer.Sprint(111.223))
		})

		t.Run("plural", func(t *testing.T) {
			t.Parallel()

			t.Run("russian", func(t *testing.T) {
				t.Parallel()

				printer := GetPrinter(language.Russian)
				assert.Equal(t, "всего 100 пауков", printer.Sprintf("test plural", 100))
				assert.Equal(t, "нет пауков", printer.Sprintf("test plural", 0))
				assert.Equal(t, "паучок", printer.Sprintf("test plural", 1))
				assert.Equal(t, "всего пара пауков", printer.Sprintf("test plural", 2))
			})

			t.Run("english", func(t *testing.T) {
				t.Parallel()

				printer := GetPrinter(language.English)
				assert.Equal(t, "exactly 100 spiders", printer.Sprintf("test plural", 100))
				assert.Equal(t, "no spiders", printer.Sprintf("test plural", 0))
				assert.Equal(t, "spider", printer.Sprintf("test plural", 1))
				assert.Equal(t, "just pair of spiders", printer.Sprintf("test plural", 2))
			})
		})
	})

	t.Run("languages ok", func(t *testing.T) {
		t.Parallel()

		languages := GetLanguages()
		require.Len(t, languages, len(supportedLanguages))
		require.ElementsMatch(t, supportedLanguages, languages)
	})
}
