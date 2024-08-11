package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func init() { RegisterConfigType(GotifyConfig{}) }

var (
	_ NotificationConfig  = GotifyConfig{}
	_ NotificationService = GotifyService{}
)

type GotifyConfig struct {
	BaseURL  string         `json:"base_url"`
	Token    string         `json:"token"`
	Priority int            `json:"priority,omitempty"`
	Extras   map[string]any `json:"extras,omitempty"`
}

func (GotifyConfig) ServiceName() string   { return "gotify" }
func (GotifyConfig) isNotificationConfig() {}

type GotifyService struct {
	Client *http.Client
}

type gotifyNotification struct {
	Title    string         `json:"title,omitempty"`
	Message  string         `json:"message"`
	Priority int            `json:"priority,omitempty"`
	Extras   map[string]any `json:"extras,omitempty"`
}

func (s GotifyService) ServiceName() string { return "gotify" }
func (s GotifyService) SendNotification(ctx context.Context, n Notification, c NotificationConfig) error {
	config, ok := c.(GotifyConfig)
	if !ok {
		return ErrInvalidConfig
	}

	u, err := url.Parse(config.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse endpoint: %w", err)
	}

	u.Path += "/message"

	q := u.Query()
	q.Set("token", config.Token)
	u.RawQuery = q.Encode()

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

	r, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode > 299 {
		return consumeHTTPUnknownStatusError(r)
	}

	return nil
}
