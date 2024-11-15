package notification

import (
	"context"
	"fmt"
	"slices"

	"go.uber.org/fx"
	"libdb.so/e2clicker/services/notification/openapi"
	"libdb.so/e2clicker/services/user"
)

// UserPreferences is the preferences of a user.
// It is JSON-serializable.
type UserPreferences struct {
	NotificationConfigs NotificationConfigs     `json:"notificationConfigs"`
	CustomNotifications map[MessageType]Message `json:"customNotifications,omitempty"`
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
}

// UserNotificationServiceConfig is the configuration for the user notification
// service.
type UserNotificationServiceConfig struct {
	fx.In

	UserNotificationStorage
	*NotificationService
	*user.UserService
}

// NewUserNotificationService creates a new user notification service.
func NewUserNotificationService(s UserNotificationServiceConfig) *UserNotificationService {
	return &UserNotificationService{
		userNotifications: s.UserNotificationStorage,
		users:             s.UserService,
		notification:      s.NotificationService,
	}
}

// NotifyUser sends a notification to a user.
func (s *UserNotificationService) NotifyUser(ctx context.Context, secret user.Secret, t MessageType) error {
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
	if custom, ok := prefs.CustomNotifications[t]; ok {
		n.Message = custom
	} else {
		n.Message, err = LoadNotification(ctx, t)
		if err != nil {
			return err
		}
	}

	return s.notification.Notify(ctx, n, prefs.NotificationConfigs)
}

// WebPushInfo returns the web push information of the server.
func (s *UserNotificationService) WebPushInfo(ctx context.Context) (openapi.PushInfo, error) {
	if s.notification.services.WebPush == nil {
		return openapi.PushInfo{}, UnknownServiceError{Service: "webPush"}
	}
	return openapi.PushInfo{
		ApplicationServerKey: s.notification.services.WebPush.VAPIDPublicKey(),
	}, nil
}

// UserPreferences returns the preferences of a user.
func (s *UserNotificationService) UserPreferences(ctx context.Context, secret user.Secret) (UserPreferences, error) {
	return s.userNotifications.UserPreferences(ctx, secret)
}

// SubscribeWebPush sets the web push subscription of a user.
// The particular subscription is identified by the device ID.
func (s *UserNotificationService) SubscribeWebPush(ctx context.Context, secret user.Secret, subscription openapi.PushSubscription) error {
	return s.userNotifications.SetUserPreferencesTx(ctx, secret, func(p *UserPreferences) error {
		ix := slices.IndexFunc(p.NotificationConfigs.WebPush,
			func(c WebPushNotificationConfig) bool {
				return c.Subscription.DeviceID == subscription.DeviceID
			},
		)
		if ix != -1 {
			p.NotificationConfigs.WebPush[ix].Subscription = subscription
		} else {
			p.NotificationConfigs.WebPush = append(p.NotificationConfigs.WebPush, WebPushNotificationConfig{
				Subscription: subscription,
			})
		}
		return nil
	})
}

// UnsubscribeWebPush removes the web push subscription of a user.
func (s *UserNotificationService) UnsubscribeWebPush(ctx context.Context, secret user.Secret, deviceID string) error {
	return s.userNotifications.SetUserPreferencesTx(ctx, secret, func(p *UserPreferences) error {
		p.NotificationConfigs.WebPush = slices.DeleteFunc(p.NotificationConfigs.WebPush,
			func(c WebPushNotificationConfig) bool {
				return c.Subscription.DeviceID == deviceID
			},
		)
		return nil
	})
}
