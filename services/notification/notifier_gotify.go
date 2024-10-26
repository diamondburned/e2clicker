package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"unique"

	"libdb.so/e2clicker/internal/validating"
)

var gotifyServiceName = registerNotificationConfig[*GotifyNotificationConfig]("gotify", nil)

type GotifyNotificationConfig struct {
	BaseURL  string         `json:"base_url"`
	Token    string         `json:"token"`
	Priority int            `json:"priority,omitempty"`
	Extras   map[string]any `json:"extras,omitempty"`
}

var _ validating.Validator = (*GotifyNotificationConfig)(nil)

func (*GotifyNotificationConfig) isNotificationConfig() {}

func (c *GotifyNotificationConfig) ServiceName() unique.Handle[string] {
	return gotifyServiceName
}

func (c *GotifyNotificationConfig) MarshalJSON() ([]byte, error) {
	return NotificationConfigJSON{gotifyServiceName, c}.MarshalJSON()
}

func (c *GotifyNotificationConfig) Validate() error {
	if _, err := url.Parse(c.BaseURL); err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}
	return nil
}

type GotifyService struct {
	http *http.Client `do:""`
}

func (s *GotifyService) Notify(ctx context.Context, n *Notification, config *GotifyNotificationConfig) error {
	if err := config.Validate(); err != nil {
		return ConfigError{err: err}
	}

	u, err := url.Parse(config.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse endpoint: %w", err)
	}

	u.Path += "/message"

	q := u.Query()
	q.Set("token", config.Token)
	u.RawQuery = q.Encode()

	type gotifyNotification struct {
		Title    string         `json:"title,omitempty"`
		Message  string         `json:"message"`
		Priority int            `json:"priority,omitempty"`
		Extras   map[string]any `json:"extras,omitempty"`
	}

	b, err := json.Marshal(gotifyNotification{
		Title:    n.Title,
		Message:  n.Message,
		Priority: config.Priority,
		Extras:   config.Extras,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	r, err := s.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode > 299 {
		return consumeHTTPUnknownStatusError(r)
	}

	return nil
}
