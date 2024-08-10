package notification

import (
	"context"
	"net/http"
)

func init() {
	RegisterService(GotifyConfig{}, (*GotifyService)(nil))
}

type GotifyConfig struct {
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
	Priority int    `json:"priority"`
}

func (s GotifyConfig) ServiceName() string   { return "gotify" }
func (s GotifyConfig) isNotificationConfig() {}

type GotifyService struct {
	Client *http.Client
}

func (s *GotifyService) ServiceName() string { return "gotify" }
func (s *GotifyService) SendNotification(ctx context.Context, n Notification, c NotificationConfig) error {
	config, ok := c.(GotifyConfig)
	if !ok {
		return ErrInvalidConfig
	}

	_ = config
	panic("not implemented")
}
