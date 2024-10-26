package notification

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"libdb.so/e2clicker/internal/publicerrors"
)

func init() {
	publicerrors.MarkValuesPublic(ErrUnknownService)
	publicerrors.MarkTypePublic[HTTPUnknownStatusError]()
	publicerrors.MarkTypePublic[ConfigError]()
}

// ErrUnknownService is returned when an unknown service is requested.
var ErrUnknownService = errors.New("unknown service")

// ConfigError is returned when a notification service is given an invalid
// configuration or the configuration fails validation.
type ConfigError struct {
	err error
}

func (e ConfigError) Error() string {
	s := "invalid config"
	if e.err != nil {
		s += ": " + e.err.Error()
	}
	return s
}

func (e ConfigError) Unwrap() error {
	return e.err
}

// HTTPUnknownStatusError is returned when an unknown HTTP status code is
// returned by an API.
type HTTPUnknownStatusError struct {
	// StatusCode is the HTTP status code of the API response.
	StatusCode int `json:"statusCode"`
	// Body is the body of the API response.
	// It is truncated to [HTTPErrorMaxBodySize] bytes.
	Body string `json:"body"`
}

const HTTPErrorMaxBodySize = 1024

func (e HTTPUnknownStatusError) Error() string {
	return fmt.Sprintf("unknown HTTP status code %d returned: %s", e.StatusCode, e.Body)
}

func consumeHTTPUnknownStatusError(resp *http.Response) error {
	limReader := io.LimitReader(resp.Body, HTTPErrorMaxBodySize)
	body, _ := io.ReadAll(limReader)
	return HTTPUnknownStatusError{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}
}
