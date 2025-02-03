package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"e2clicker.app/internal/validating"
)

type PushoverNotificationConfig struct {
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Token    string `json:"token"`
	Priority int    `json:"priority,omitempty"`
	Sound    string `json:"sound,omitempty"`
	Device   string `json:"device,omitempty"`
}

var _ validating.Validator = (*PushoverNotificationConfig)(nil)

// Validate checks that the configuration is valid.
func (c *PushoverNotificationConfig) Validate() error {
	if _, err := url.Parse(c.Endpoint); err != nil {
		return fmt.Errorf("invalid endpoint: %w", err)
	}
	return nil
}

// PushoverService is a service for sending notifications via Pushover.
type PushoverService struct {
	http *http.Client `do:""`
}

// NewPushoverService creates a new Pushover service.
func NewPushoverService(c *http.Client) *PushoverService {
	return &PushoverService{http: c}
}

func (s PushoverService) Notify(ctx context.Context, n Notification, config PushoverNotificationConfig) error {
	if err := config.Validate(); err != nil {
		return ConfigError{err: err}
	}

	type pushoverNotification struct {
		User     string `json:"user"`
		Token    string `json:"token"`
		Message  string `json:"message"`
		Title    string `json:"title,omitempty"`
		Priority int    `json:"priority,omitempty"`
		Sound    string `json:"sound,omitempty"`
		Device   string `json:"device,omitempty"`
	}

	b, err := json.Marshal(pushoverNotification{
		Title:    n.Message.Title,
		Message:  n.Message.Message,
		User:     config.User,
		Token:    config.Token,
		Priority: config.Priority,
		Sound:    config.Sound,
		Device:   config.Device,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.Endpoint, bytes.NewReader(b))
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
