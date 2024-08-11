package notification

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ErrInvalidConfig is returned when a notification service is given an
// invalid configuration.
var ErrInvalidConfig = errors.New("invalid config")

// UnknownServiceError is returned when a service is unknown.
type UnknownServiceError struct {
	ServiceName string
}

func (e UnknownServiceError) Error() string {
	return "unknown service: " + e.ServiceName
}

// HTTPUnknownStatusError is returned when an unknown HTTP status code is
// returned by an API.
type HTTPUnknownStatusError struct {
	// StatusCode is the HTTP status code of the API response.
	StatusCode int
	// Body is the body of the API response.
	// It is truncated to [HTTPErrorMaxBodySize] bytes.
	Body []byte
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
		Body:       body,
	}
}
