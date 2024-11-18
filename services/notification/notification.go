package notification

import (
	"context"

	"libdb.so/e2clicker/services/notification/openapi"
)

// Notification describes a notification message to be sent to the user.
type Notification = openapi.Notification

// LoadNotification loads a notification message of the given type.
func LoadNotification(ctx context.Context, t openapi.NotificationType) (openapi.NotificationMessage, error) {
	switch t {
	case openapi.WelcomeMessage:
		return openapi.NotificationMessage{
			Title:   "Welcome! 😄🌈❤️",
			Message: "e2clicker can send you notifications to remind you now!",
		}, nil
	case openapi.ReminderMessage:
		return openapi.NotificationMessage{
			Title:   "Reminder!",
			Message: "Don't forget to take your hormone dose!",
		}, nil
	case openapi.AccountNoticeMessage:
		return openapi.NotificationMessage{
			Title:   "Account Notice",
			Message: "Please check your e2clicker account.",
		}, nil
	case openapi.WebPushExpiringMessage:
		return openapi.NotificationMessage{
			Title:   "Notifications will stop working soon 😟",
			Message: "Your browser's push subscription is expiring soon. You need to refresh it!",
		}, nil
	default:
		return openapi.NotificationMessage{}, ErrUnknownNotificationType
	}
}
