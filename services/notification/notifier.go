package notification

import (
	"context"
	"log/slog"
	"reflect"

	"go.uber.org/fx"
	"libdb.so/e2clicker/internal/validating"
)

// Notifier sends notifications to users.
type Notifier interface {
	// Notify sends a notification using this service.
	Notify(context.Context, *Notification, NotificationConfig) error
}

// NotificationService is a collection of NotificationServices.
// It implements the [Notifier] interface.
type NotificationService struct {
	logger   *slog.Logger
	gotify   *GotifyService
	pushover *PushoverService
}

// NotificationServiceConfig is the configuration for the notification service.
type NotificationServiceConfig struct {
	fx.In
	Gotify   *GotifyService   `optional:"true"`
	Pushover *PushoverService `optional:"true"`
}

var _ Notifier = (*NotificationService)(nil)

// NewNotificationService creates a new notification service.
func NewNotificationService(c NotificationServiceConfig, logger *slog.Logger) *NotificationService {
	return &NotificationService{
		logger:   logger,
		gotify:   c.Gotify,
		pushover: c.Pushover,
	}
}

// Notify implements [Notifier].
func (m NotificationService) Notify(ctx context.Context, n *Notification, c NotificationConfig) error {
	if err := validating.Validate(ctx, c); err != nil {
		return ConfigError{err}
	}

	switch c := c.(type) {
	case *GotifyNotificationConfig:
		if m.gotify != nil {
			return m.gotify.Notify(ctx, n, c)
		}
	case *PushoverNotificationConfig:
		if m.pushover != nil {
			return m.pushover.Notify(ctx, n, c)
		}
	}

	m.logger.WarnContext(ctx,
		"BUG: unknown notification service config",
		"gotify", m.gotify != nil,
		"pushover", m.pushover != nil,
		"config.type", reflect.TypeOf(c))

	return ErrUnknownService
}
