package i18n_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	. "github.com/derfenix/goi18n"
	"github.com/derfenix/goi18n/internal"
)

func TestTryTranslateError(t *testing.T) {
	t.Parallel()

	require.NoError(t, Init(internal.TestFS))

	err := NewError("test").WithParams("book")
	wrapped := errors.Wrap(err, "foo bar")

	errorString, translated := TryTranslateError(context.Background(), wrapped)
	require.True(t, translated)
	assert.Equal(t, "Тест book", errorString)

	errorString, translated = TryTranslateError(ContextWithLang(context.Background(), language.English), wrapped)
	require.True(t, translated)
	assert.Equal(t, "Test of the book", errorString)
}
