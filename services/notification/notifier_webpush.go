package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"libdb.so/e2clicker/services/notification/openapi"

	e2clickermodule "libdb.so/e2clicker/nix/modules/e2clicker"
)

// WebPushNotificationConfig is a configuration for the Push API service.
type WebPushNotificationConfig struct {
	Subscription openapi.PushSubscription `json:"subscription"`
}

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
	http http.Client
	keys e2clickermodule.WebPushKeysSubmodule
}

// NewWebPushSevice creates a new Push API service.;
func NewWebPushSevice(config e2clickermodule.Notification) (*WebPushService, error) {
	timeout, err := time.ParseDuration(config.ClientTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid client timeout %q: %w", config.ClientTimeout, err)
	}

	var keys e2clickermodule.WebPushKeysSubmodule

	switch value := config.WebPushKeys.Value.(type) {
	case e2clickermodule.WebPushKeysSubmodule:
		keys = value
	case e2clickermodule.WebPushKeysPath:
		b, err := os.ReadFile(string(value))
		if err != nil {
			return nil, fmt.Errorf("cannot read WebPushKeys file at %s: %w", value, err)
		}
		if err := json.Unmarshal(b, &keys); err != nil {
			return nil, fmt.Errorf("cannot unmarshal WebPushKeys at %s: %w", value, err)
		}
	default:
		panic("unreachable")
	}

	return &WebPushService{
		http: http.Client{Timeout: timeout},
		keys: keys,
	}, nil
}

func (s WebPushService) Notify(ctx context.Context, n Notification, config WebPushNotificationConfig) error {
	if config.Subscription.ExpirationTime != nil {
		expMs := *config.Subscription.ExpirationTime
		expSecInt, expSecFrac := math.Modf(expMs / 1000)
		exp := time.Unix(
			int64(expSecInt),
			int64(expSecFrac*float64(time.Second)),
		)

		if exp.Before(time.Now()) {
			return &WebPushSubscriptionExpired{exp}
		}
	}

	m, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("cannot marshal notification: %w", err)
	}

	resp, err := webpush.SendNotificationWithContext(
		ctx, m,
		convertSubscription(config.Subscription),
		&webpush.Options{
			HTTPClient:      &s.http,
			Urgency:         webpush.UrgencyHigh,
			VAPIDPublicKey:  s.keys.PublicKey,
			VAPIDPrivateKey: s.keys.PrivateKey,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot send notification: %w", err)
	}

	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return consumeHTTPUnknownStatusError(resp)
	}

	return nil
}
