package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func init() { RegisterConfigType(PushoverConfig{}) }

var (
	_ NotificationConfig  = PushoverConfig{}
	_ NotificationService = PushoverService{}
)

type PushoverConfig struct {
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Token    string `json:"token"`
	Priority int    `json:"priority,omitempty"`
	Sound    string `json:"sound,omitempty"`
	Device   string `json:"device,omitempty"`
}

func (PushoverConfig) ServiceName() string   { return "pushover" }
func (PushoverConfig) isNotificationConfig() {}

type PushoverService struct {
	Client *http.Client
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

func (s PushoverService) ServiceName() string { return "pushover" }
func (s PushoverService) SendNotification(ctx context.Context, n Notification, c NotificationConfig) error {
	config, ok := c.(PushoverConfig)
	if !ok {
		return ErrInvalidConfig
	}

	b, err := json.Marshal(pushoverNotification{
		Title:    n.Title,
		Message:  n.Message,
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
