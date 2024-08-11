// Package notification provides a way to send notifications to users.
package notification

import (
	"context"
	"reflect"
)

// NotificationService sends notifications to users.
type NotificationService interface {
	// SendNotification sends a notification using this service.
	SendNotification(context.Context, Notification, NotificationConfig) error
}

// Notification is a message to be sent to a user.
type Notification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// NotificationServiceMux is a collection of NotificationServices.
// It implements the [NotificationService] interface.
type NotificationServiceMux struct {
	m map[reflect.Type]NotificationService
}

// MuxableNotificationService extends NotificationService to include the
// ServiceName method.
type MuxableNotificationService interface {
	NotificationService
	// ServiceName returns the name of the service.
	ServiceName() string
}

// NewNotificationServiceMux creates a new NotificationServiceMux with the given services.
// If the mux is given an unknown service, it will panic.
func NewNotificationServiceMux(services ...MuxableNotificationService) NotificationServiceMux {
	m := map[reflect.Type]NotificationService{}
	for _, service := range services {
		configType, ok := knownConfigTypes[service.ServiceName()]
		if !ok {
			panic("unknown service: " + service.ServiceName())
		}
		m[configType] = service
	}
	return NotificationServiceMux{m: m}
}

// SendNotification implements [NotificationService].
func (m NotificationServiceMux) SendNotification(ctx context.Context, n Notification, c NotificationConfig) error {
	ct := reflect.TypeOf(c)
	cs, ok := m.m[ct]
	if !ok {
		return UnknownServiceError{ServiceName: c.ServiceName()}
	}
	return cs.SendNotification(ctx, n, c)
}
