package notification

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"libdb.so/e2clicker/internal/publicerrors"
)

func init() {
	publicerrors.MarkTypePublic[UnknownServiceError]()
	publicerrors.MarkTypePublic[HTTPUnknownStatusError]()
	publicerrors.MarkTypePublic[ConfigError]()
	publicerrors.MarkTypePublic[WebPushSubscriptionExpired]()
}

// UnknownServiceError is returned when an unknown service is requested.
type UnknownServiceError struct {
	Service string `json:"service"`
}

func (e UnknownServiceError) Error() string {
	return fmt.Sprintf("unknown service %q", e.Service)
}

// ConfigError is returned when a notification service is given an invalid
// configuration or the configuration fails validation.
type ConfigError struct {
	Service string `json:"service"`
	err     error
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

// WebPushSubscriptionExpired is returned when a WebPush subscription has
// expired.
type WebPushSubscriptionExpired struct {
	ExpiredAt time.Time `json:"expiredAt"`
}

func (e WebPushSubscriptionExpired) Error() string {
	return fmt.Sprintf("push subscription expired at %s", e.ExpiredAt.Format(time.RFC3339))
}
