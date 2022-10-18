package i18n

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type TranslatableError interface {
	error
	Translatable
}

func TryTranslateError(ctx context.Context, err error) (string, bool) {
	var translatable TranslatableError

	if errors.As(err, &translatable) {
		return translatable.Translate(ctx), true
	}

	return "", false
}

func NewError(key string) *Error {
	return &Error{key: key}
}

type Error struct {
	key    string
	params []interface{}
}

func (e *Error) Error() string {
	return e.key
}

func (e *Error) WithParams(params ...interface{}) *Error {
	newErr := *e

	newErr.params = params

	return &newErr
}

func (e *Error) Translate(ctx context.Context) string {
	printer := PrinterFromContext(ctx)

	translatedParams := make([]interface{}, len(e.params))

	for idx, param := range e.params {
		switch typed := param.(type) {
		case string:
			translatedParams[idx] = printer.Sprintf(typed)
		case fmt.Stringer:
			translatedParams[idx] = printer.Sprintf(typed.String())
		default:
			translatedParams[idx] = param
		}
	}

	return printer.Sprintf(e.key, translatedParams...)
}
