package notification

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"libdb.so/e2clicker/services/api/openapi"
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
	// SetUserPreferences sets the preferences of a user.
	SetUserPreferences(ctx context.Context, userSecret user.Secret, prefs UserPreferences) error
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

// UserPreferences returns the preferences of a user.
func (s *UserNotificationService) UserPreferences(ctx context.Context, secret user.Secret) (UserPreferences, error) {
	return s.userNotifications.UserPreferences(ctx, secret)
}

// AddWebPushSubscription adds a web push subscription to a user.
func (s *UserNotificationService) AddWebPushSubscription(ctx context.Context, secret user.Secret, subscription openapi.PushSubscription) error {
	panic("implement me")

	// p, err := s.UserPreferences(ctx, secret)
	// if err != nil {
	// 	return err
	// }
	//
	// var found bool
	// for i := range p.NotificationConfigs.WebPush {
	// 	c := &p.NotificationConfigs.WebPush[i]
	// 	if c.Subscription.Endpoint == subscription.Endpoint {
	// 		found = true
	// 		c.Subscription = subscription
	// 		break
	// 	}
	// }
	// if !found {
	// 	p.NotificationConfigs.WebPush = append(p.NotificationConfigs.WebPush, WebPushNotificationConfig{
	// 		Subscription: subscription,
	// 	})
	// }
	//
	// if oldSubscription == nil {
	// 	oldSubscription = &openapi.PushSubscription{}
	// 	return nil
	// }
	//
	// return s.userNotifications.SetUserPreferences(ctx, secret, prefs)
}
