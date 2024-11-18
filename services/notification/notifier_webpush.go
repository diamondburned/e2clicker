package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"libdb.so/e2clicker/services/notification/openapi"

	e2clickermodule "libdb.so/e2clicker/nix/modules/e2clicker"
)

// WebPushNotificationConfig is a configuration for the Push API service.
type WebPushNotificationConfig = openapi.PushSubscription

func convertSubscription(subscription openapi.PushSubscription) *webpush.Subscription {
	return &webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			P256dh: subscription.Keys.P256Dh,
			Auth:   subscription.Keys.Auth,
		},
	}
}

// WebPushService is a service for sending notifications via the Push API.
type WebPushService struct {
	http *http.Client
	keys *e2clickermodule.VAPIDKeysSubmodule
}

// NewWebPushSevice creates a new Push API service.
func NewWebPushSevice(config e2clickermodule.Notification) (*WebPushService, error) {
	if !config.WebPush.Enable {
		return nil, nil
	}

	timeout, err := time.ParseDuration(config.ClientTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid client timeout %q: %w", config.ClientTimeout, err)
	}

	var keys *e2clickermodule.VAPIDKeysSubmodule

	switch value := config.WebPush.VAPIDKeys.Value.(type) {
	case e2clickermodule.VAPIDKeysSubmodule:
		keys = &value
	case e2clickermodule.VAPIDKeysPath:
		b, err := os.ReadFile(string(value))
		if err != nil {
			return nil, fmt.Errorf("cannot read WebPushKeys file at %s: %w", value, err)
		}
		keys = new(e2clickermodule.VAPIDKeysSubmodule)
		if err := json.Unmarshal(b, keys); err != nil {
			return nil, fmt.Errorf("cannot unmarshal WebPushKeys at %s: %w", value, err)
		}
	default:
		panic("unreachable")
	}

	return &WebPushService{
		http: &http.Client{Timeout: timeout},
		keys: keys,
	}, nil
}

// VAPIDPublicKey returns the VAPID public key.
func (s WebPushService) VAPIDPublicKey() string {
	return s.keys.PublicKey
}

func (s WebPushService) Notify(ctx context.Context, n Notification, config WebPushNotificationConfig) error {
	if !config.ExpirationTime.IsZero() && config.ExpirationTime.Before(time.Now()) {
		return &WebPushSubscriptionExpired{config.ExpirationTime}
	}

	m, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("cannot marshal notification: %w", err)
	}

	opts := &webpush.Options{
		HTTPClient:      s.http,
		Urgency:         webpush.UrgencyHigh,
		Subscriber:      n.Username,
		VAPIDPublicKey:  s.keys.PublicKey,
		VAPIDPrivateKey: s.keys.PrivateKey,
	}

	resp, err := webpush.SendNotificationWithContext(ctx, m, convertSubscription(config), opts)
	if err != nil {
		return fmt.Errorf("cannot send notification: %w", err)
	}

	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return consumeHTTPUnknownStatusError(resp)
	}

	return nil
}
