package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"e2clicker.app/services/notification/openapi"
	"e2clicker.app/services/user"
	"go.uber.org/fx"
)

// UserPreferences is the preferences of a user.
// It is JSON-serializable.
type UserPreferences struct {
	NotificationConfigs NotificationConfigs         `json:"notificationConfigs"`
	CustomNotifications openapi.CustomNotifications `json:"customNotifications,omitempty"`
}

type UserNotificationStorage interface {
	// UserPreferences returns the preferences of a user.
	UserPreferences(ctx context.Context, userSecret user.Secret) (UserPreferences, error)
	// SetUserPreferencesTx sets the preferences of a user inside a transaction.
	// The function set is called with the current preferences and should modify
	// the given preferences, all within the same transaction.
	SetUserPreferencesTx(ctx context.Context, userSecret user.Secret, set func(*UserPreferences) error) error
}

// UserNotificationService is a service that sends notifications to users.
type UserNotificationService struct {
	userNotifications UserNotificationStorage
	users             *user.UserService
	notification      *NotificationService
	logger            *slog.Logger
}

// UserNotificationServiceConfig is the configuration for the user notification
// service.
type UserNotificationServiceConfig struct {
	fx.In

	UserNotificationStorage
	*NotificationService
	*user.UserService
	*slog.Logger
}

// NewUserNotificationService creates a new user notification service.
func NewUserNotificationService(s UserNotificationServiceConfig) *UserNotificationService {
	return &UserNotificationService{
		userNotifications: s.UserNotificationStorage,
		users:             s.UserService,
		notification:      s.NotificationService,
		logger:            s.Logger,
	}
}

// NotifyUser sends a notification to a user.
func (s *UserNotificationService) NotifyUser(ctx context.Context, secret user.Secret, t openapi.NotificationType) error {
	prefs, err := s.userNotifications.UserPreferences(ctx, secret)
	if err != nil {
		return err
	}

	if prefs.NotificationConfigs.IsEmpty() {
		return nil
	}

	u, err := s.users.User(ctx, secret)
	if err != nil {
		return fmt.Errorf("failed to get user for notification: %w", err)
	}

	n := Notification{
		Type:     t,
		Username: u.Name,
	}
	if custom, ok := prefs.CustomNotifications[string(t)]; ok {
		n.Message = custom
	} else {
		n.Message, err = LoadNotification(ctx, t)
		if err != nil {
			return err
		}
	}

	return s.notification.Notify(ctx, n, prefs.NotificationConfigs)
}

// UserPreferences returns the preferences of a user.
func (s *UserNotificationService) UserPreferences(ctx context.Context, secret user.Secret) (UserPreferences, error) {
	return s.userNotifications.UserPreferences(ctx, secret)
}

// SetUserPreferences sets the preferences of a user.
func (s *UserNotificationService) SetUserPreferences(ctx context.Context, secret user.Secret, preferences *UserPreferences) error {
	return s.SetUserPreferencesSafe(ctx, secret, preferences, nil)
}

// SetUserPreferencesSafe sets the preferences of a user in a safer manner than
// [SetUserPreferences]. It requires the old preferences to be passed in order
// to prevent overwriting changes made by other clients.
//
// Realistically, this doesn't happen unless the user is deliberately trying to
// cause the issue.
func (s *UserNotificationService) SetUserPreferencesSafe(ctx context.Context, secret user.Secret, newPreferences, oldPreferences *UserPreferences) error {
	return s.userNotifications.SetUserPreferencesTx(ctx, secret, func(p *UserPreferences) error {
		if oldPreferences != nil {
			b1, _ := json.Marshal(oldPreferences)
			b2, _ := json.Marshal(p)
			if !bytes.Equal(b1, b2) {
				return fmt.Errorf("preferences have changed since you last read them")
			}
		}
		*p = *newPreferences
		return nil
	})
}

// WebPushInfo returns the web push information of the server.
func (s *UserNotificationService) WebPushInfo(ctx context.Context) (openapi.PushInfo, error) {
	if s.notification.services.WebPush == nil {
		return openapi.PushInfo{}, ErrWebPushNotAvailable
	}
	return openapi.PushInfo{
		ApplicationServerKey: s.notification.services.WebPush.VAPIDPublicKey(),
	}, nil
}

/*
// SubscribeWebPush sets the web push subscription of a user.
// The particular subscription is identified by the device ID.
func (s *UserNotificationService) SubscribeWebPush(ctx context.Context, secret user.Secret, subscription openapi.PushSubscription) error {
	s.logger.Debug(
		"Updating Web Push subscription",
		"endpoint", subscription.Endpoint,
		"expirationTime", subscription.ExpirationTime)

	return s.userNotifications.SetUserPreferencesTx(ctx, secret, func(p *UserPreferences) error {
		ix := slices.IndexFunc(p.NotificationConfigs.WebPush,
			func(c WebPushNotificationConfig) bool {
				return c.DeviceID == subscription.DeviceID
			},
		)
		if ix != -1 {
			p.NotificationConfigs.WebPush[ix] = subscription
		} else {
			p.NotificationConfigs.WebPush = append(p.NotificationConfigs.WebPush, subscription)
		}
		return nil
	})
}

// UnsubscribeWebPush removes the web push subscription of a user.
func (s *UserNotificationService) UnsubscribeWebPush(ctx context.Context, secret user.Secret, deviceID string) error {
	return s.userNotifications.SetUserPreferencesTx(ctx, secret, func(p *UserPreferences) error {
		p.NotificationConfigs.WebPush = slices.DeleteFunc(p.NotificationConfigs.WebPush,
			func(c WebPushNotificationConfig) bool { return c.DeviceID == deviceID },
		)
		return nil
	})
}
*/
