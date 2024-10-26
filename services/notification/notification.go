package notification

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unique"
)

// Notification is a message to be sent to a user.
type Notification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// NotificationConfig is a configuration for a notification.
type NotificationConfig interface {
	json.Marshaler
	isNotificationConfig()
}

var configFactories = map[unique.Handle[string]]func() NotificationConfig{}

func registerNotificationConfig[T NotificationConfig](name string, newFn func() NotificationConfig) unique.Handle[string] {
	if newFn == nil {
		rt := reflect.TypeFor[T]()
		if rt.Kind() != reflect.Ptr {
			panic("notification config must be a pointer")
		}
		rt = rt.Elem()
		newFn = func() NotificationConfig {
			rv := reflect.New(rt)
			return rv.Interface().(NotificationConfig)
		}
	}

	nameHandle := unique.Make(name)
	configFactories[nameHandle] = newFn

	return nameHandle
}

// NotificationConfigJSON is a JSON representation of a [NotificationConfig].
// When handling JSON, this type must be used for marshaling to succeed.
type NotificationConfigJSON struct {
	ServiceName unique.Handle[string]
	NotificationConfig
}

func (n NotificationConfigJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Service string             `json:"service"`
		Config  NotificationConfig `json:"config"`
	}{
		Service: n.ServiceName.Value(),
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

	newCfg, ok := configFactories[unique.Make(raw.Service)]
	if !ok {
		return fmt.Errorf("unknown notification service %q", raw.Service)
	}

	config := newCfg()
	if err := json.Unmarshal(raw.Config, config); err != nil {
		return err
	}

	n.NotificationConfig = config
	return nil
}
