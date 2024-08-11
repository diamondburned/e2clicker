package notification

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var knownConfigTypes = map[string]reflect.Type{}

// RegisterConfigType adds the given config type to the list of known
// configuration types for unmarshaling. The service name is taken from the
// config type for actual use.
func RegisterConfigType(config NotificationConfig) {
	serviceName := config.ServiceName()
	if serviceName == "" {
		panic("service name must not be empty")
	}
	knownConfigTypes[serviceName] = reflect.TypeOf(config)
}

// NotificationConfigJSON is a JSON representation of a NotificationConfig.
// When handling JSON, this type must be used for marshaling to succeed.
type NotificationConfigJSON struct {
	NotificationConfig
}

// IsValid returns true if the NotificationConfig is not nil and is known.
func (n NotificationConfigJSON) IsValid() bool {
	return n.NotificationConfig != nil && knownConfigTypes[n.NotificationConfig.ServiceName()] != nil
}

func (n NotificationConfigJSON) isNotificationConfig() {}

func (n NotificationConfigJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Service string             `json:"service"`
		Config  NotificationConfig `json:"config"`
	}{
		Service: n.NotificationConfig.ServiceName(),
		Config:  n.NotificationConfig,
	})
}

func (n *NotificationConfigJSON) UnmarshalJSON(data []byte) error {
	var raw struct {
		Service string          `json:"service"`
		Config  json.RawMessage `json:"config"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	rt, ok := knownConfigTypes[raw.Service]
	if !ok {
		return fmt.Errorf("unknown notification service %q", raw.Service)
	}

	var wasPtr bool
	if rt.Kind() == reflect.Ptr {
		wasPtr = true
		rt = rt.Elem()
	}

	config := reflect.New(rt)
	if err := json.Unmarshal(raw.Config, config.Interface()); err != nil {
		return err
	}

	if wasPtr {
		config = config.Addr()
	}

	// Bug check.
	if config.Type() != rt {
		panic("BUG: config type mismatch")
	}

	n.NotificationConfig = config.Interface().(NotificationConfig)
	return nil
}

// NotificationConfig is a configuration for a notification.
type NotificationConfig interface {
	// ServiceName returns the name of the notification service that this
	// configuration is for.
	ServiceName() string

	isNotificationConfig()
}
