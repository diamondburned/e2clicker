package notification

import (
	"context"
	"encoding/json"
	"errors"

	"libdb.so/e2clicker/services/notification/openapi"
)

// Notification describes a notification message to be sent to the user.
type Notification struct {
	// Type is the type of the notification message.
	Type MessageType
	// Message is the notification message.
	Message Message
	// Username is the name of the user that this notification is for.
	Username string
}

// MarshalJSON marshals the notification to JSON according to the OpenAPI
// schema.
func (n Notification) MarshalJSON() ([]byte, error) {
	return json.Marshal(openapi.Notification{
		Type: openapi.NotificationType(n.Type),
		Message: openapi.NotificationMessage{
			Title:   n.Message.Title,
			Message: n.Message.Message,
		},
		Username: n.Username,
	})
}

// UnmarshalJSON unmarshals the notification from JSON according to the OpenAPI
// schema.
func (n *Notification) UnmarshalJSON(data []byte) error {
	var openapiNotification openapi.Notification
	if err := json.Unmarshal(data, &openapiNotification); err != nil {
		return err
	}

	*n = Notification{
		Type: MessageType(openapiNotification.Type),
		Message: Message{
			Title:   openapiNotification.Message.Title,
			Message: openapiNotification.Message.Message,
		},
		Username: openapiNotification.Username,
	}
	return nil
}

// MessageType is the type of the notification message.
type MessageType string

const (
	// WelcomeMessage is sent to welcome the user.
	// Realistically, it is used as a test message.
	WelcomeMessage MessageType = "welcome"
	// ReminderMessage is sent to remind the user of their hormone dose.
	ReminderMessage MessageType = "reminder"
	// AccountNoticeMessage is sent to notify the user that they need to
	// check their account.
	AccountNoticeMessage MessageType = "account_notice"
	// WebPushExpiringMessage is sent to notify the user that their web push
	// subscription is expiring.
	WebPushExpiringMessage MessageType = "web_push_expiring"
)

// ErrUnknownNotificationType is returned when the notification type is unknown.
var ErrUnknownNotificationType = errors.New("unknown notification type")

// Message is the message of the notification to be sent to a user.
type Message struct {
	Title   string
	Message string
}

// LoadNotification loads a notification message of the given type.
func LoadNotification(ctx context.Context, t MessageType) (Message, error) {
	switch t {
	case WelcomeMessage:
		return Message{
			Title:   "Welcome! 😄🌈❤️",
			Message: "e2clicker can send you notifications to remind you now!",
		}, nil
	case ReminderMessage:
		return Message{
			Title:   "Reminder!",
			Message: "Don't forget to take your hormone dose!",
		}, nil
	case AccountNoticeMessage:
		return Message{
			Title:   "Account Notice",
			Message: "Please check your e2clicker account.",
		}, nil
	case WebPushExpiringMessage:
		return Message{
			Title:   "Notifications will stop working soon 😟",
			Message: "Your browser's push subscription is expiring soon. You need to refresh it!",
		}, nil
	default:
		return Message{}, ErrUnknownNotificationType
	}
}
