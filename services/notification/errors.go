package notification

import (
	"fmt"
	"io"
	"net/http"

	"libdb.so/hrtclicker/v2/internal/publicerrors"
)

func init() {
	publicerrors.MarkTypePublic[UnknownServiceError]()
	publicerrors.MarkTypePublic[HTTPUnknownStatusError]()
	publicerrors.MarkTypePublic[ConfigError]()
}

// ConfigError is returned when a notification service is given an invalid
// configuration or the configuration fails validation.
type ConfigError struct {
	ServiceName string `json:"serviceName"`
	err         error
}

func (e ConfigError) Error() string {
	s := fmt.Sprintf("invalid config for service %q", e.ServiceName)
	if e.err != nil {
		s += ": " + e.err.Error()
	}
	return s
}

func (e ConfigError) Unwrap() error {
	return e.err
}

// UnknownServiceError is returned when a service is unknown.
type UnknownServiceError struct {
	ServiceName string `json:"serviceName"`
}

func (e UnknownServiceError) Error() string {
	return fmt.Sprintf("unknown service: %q", e.ServiceName)
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
