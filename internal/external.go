package internal

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

var ErrInvalidResponseCode = errors.New("invalid response status code")

func NewExternalLoader(baseURL string, header http.Header) *ExternalLoader {
	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 2 * time.Second,
			}).DialContext,
			DisableKeepAlives:      true,
			MaxIdleConns:           1,
			MaxIdleConnsPerHost:    1,
			MaxConnsPerHost:        2,
			IdleConnTimeout:        30 * time.Second,
			ResponseHeaderTimeout:  5 * time.Second,
			MaxResponseHeaderBytes: 1024 * 3,
			WriteBufferSize:        100,
			ReadBufferSize:         1024 * 8,
		},
		Timeout: time.Second * 10,
	}

	return NewExternalLoaderWithClient(baseURL, header, &client)
}

func NewExternalLoaderWithClient(baseURL string, header http.Header, client *http.Client) *ExternalLoader {
	return &ExternalLoader{baseURL: baseURL, client: client, header: header}
}

type ExternalLoader struct {
	baseURL string
	header  http.Header
	client  *http.Client
}

func (e *ExternalLoader) Load(builder *catalog.Builder) error {
	languages := builder.Languages()

	for _, lang := range languages {
		if err := e.load(lang, builder); err != nil {
			return errors.WithMessagef(err, "load translation for %s", lang.String())
		}
	}

	return nil
}

func (e *ExternalLoader) load(lang language.Tag, builder *catalog.Builder) error {
	langURL, err := url.JoinPath(e.baseURL, lang.String())
	if err != nil {
		return errors.Wrap(err, "join url path")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, langURL, nil)
	if err != nil {
		return errors.Wrap(err, "new request")
	}

	if e.header != nil {
		req.Header = e.header
	}

	response, err := e.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}

	if response.StatusCode != http.StatusOK {
		return errors.Wrapf(ErrInvalidResponseCode, "got status %d", response.StatusCode)
	}

	if err := load(response.Body, lang, builder); err != nil {
		return errors.WithMessage(err, "load translation")
	}

	return nil
}
