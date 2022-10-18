//go:build i18n_extra

package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"

	"github.com/derfenix/goi18n/internal"
)

func TestOverrideTranslation(t *testing.T) {
	t.Parallel()

	SetExternalBuilder(func(builder *catalog.Builder) error {
		if err := builder.Set(language.Russian, "test", catalog.String("Тост %s")); err != nil {
			return err
		}

		return nil
	})

	require.NoError(t, Init(internal.TestFS))

	translated := GetPrinter(language.Russian).Sprintf("test", "был")
	assert.Equal(t, "Тост был", translated)
}

func TestExternalLoader(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/en":
			_, _ = w.Write([]byte(`[{"key": "Published","translation": "Foo"}]`))
		case "/ru":
			_, _ = w.Write([]byte(`[{"key": "Published","translation": "Буу"}]`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	SetExternalLoader(NewExternalLoader(srv.URL, http.Header{}))
	require.NoError(t, Init(internal.TestFS))

	{
		translated := GetPrinter(language.Russian).Sprintf("Published")
		assert.Equal(t, "Буу", translated)
	}

	{
		translated := GetPrinter(language.English).Sprintf("Published")
		assert.Equal(t, "Foo", translated)
	}

	t.Run("refresh", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.RequestURI {
			case "/en":
				_, _ = w.Write([]byte(`[{"key": "Published","translation": "Bar"}]`))
			case "/ru":
				_, _ = w.Write([]byte(`[{"key": "Published","translation": "Бяя"}]`))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))

		SetExternalLoader(NewExternalLoader(srv.URL, http.Header{}))
		require.NoError(t, RefreshTranslations())

		{
			translated := GetPrinter(language.Russian).Sprintf("Published")
			assert.Equal(t, "Бяя", translated)
		}

		{
			translated := GetPrinter(language.English).Sprintf("Published")
			assert.Equal(t, "Bar", translated)
		}
	})
}
