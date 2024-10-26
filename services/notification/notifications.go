package notification

import (
	"context"
	"errors"
)

// NotificationMessageType is the type of the notification message.
type NotificationMessageType string

const (
	// ReminderNotification is sent to remind the user of their hormone dose.
	ReminderNotification NotificationMessageType = "reminder"
	// AccountNoticeNotification is sent to notify the user that they need to
	// check their account.
	AccountNoticeNotification NotificationMessageType = "account_notice"
)

// ErrUnknownNotificationType is returned when the notification type is unknown.
var ErrUnknownNotificationType = errors.New("unknown notification type")

// LoadNotification loads a notification message of the given type.
func LoadNotification(ctx context.Context, t NotificationMessageType) (*Notification, error) {
	switch t {
	case ReminderNotification:
		return &Notification{
			Title:   "Reminder",
			Message: "Don't forget to take your hormone dose!",
		}, nil
	case AccountNoticeNotification:
		return &Notification{
			Title:   "Account Notice",
			Message: "Please check your e2clicker account.",
		}, nil
	default:
		return nil, ErrUnknownNotificationType
	}
}
