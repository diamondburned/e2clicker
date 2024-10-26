package notification

import (
	"context"
	"log/slog"
	"reflect"

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
	logger   *slog.Logger     `do:""`
	gotify   *GotifyService   `do:""`
	pushover *PushoverService `do:""`
}

var _ Notifier = (*NotificationService)(nil)

// Notify implements [Notifier].
func (m NotificationService) Notify(ctx context.Context, n *Notification, c NotificationConfig) error {
	if err := validating.Validate(ctx, c); err != nil {
		return ConfigError{err}
	}

	switch c := c.(type) {
	case *GotifyNotificationConfig:
		return m.gotify.Notify(ctx, n, c)
	case *PushoverNotificationConfig:
		return m.pushover.Notify(ctx, n, c)
	default:
		m.logger.WarnContext(ctx,
			"BUG: unknown notification service config",
			"service", reflect.TypeOf(c))
		return ErrUnknownService
	}
}
