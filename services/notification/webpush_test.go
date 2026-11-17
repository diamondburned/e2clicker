package notification

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"e2clicker.app/internal/publicerrors"
	e2clickermodule "e2clicker.app/nix/modules/e2clicker"
	"e2clicker.app/services/notification/openapi"
	"e2clicker.app/services/user"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/alecthomas/assert/v2"
	"github.com/neilotoole/slogt"
)

func testKeys(t *testing.T) (priv, pub, p256dh, auth string) {
	t.Helper()
	var err error
	priv, pub, err = webpush.GenerateVAPIDKeys()
	assert.NoError(t, err)

	// 65-byte uncompressed P-256 point (0x04 || X || Y) using valid public key bytes.
	rawPub, err := base64.RawURLEncoding.DecodeString(pub)
	assert.NoError(t, err)
	p256dh = base64.RawURLEncoding.EncodeToString(rawPub)
	auth = base64.RawURLEncoding.EncodeToString([]byte("0123456789abcdef"))
	return priv, pub, p256dh, auth
}

func TestNormalizeVAPIDSubscriber(t *testing.T) {
	sub, err := normalizeVAPIDSubscriber("mailto:admin@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "admin@example.com", sub)

	sub, err = normalizeVAPIDSubscriber("admin@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "admin@example.com", sub)

	sub, err = normalizeVAPIDSubscriber("https://e2clicker.app/contact")
	assert.NoError(t, err)
	assert.Equal(t, "https://e2clicker.app/contact", sub)

	for _, invalid := range []string{"", "Happy Sunflower", "http://e2clicker.app", "mailto:"} {
		_, err := normalizeVAPIDSubscriber(invalid)
		assert.Error(t, err)
	}
}

func TestValidatePushEndpoint(t *testing.T) {
	assert.NoError(t, ValidatePushEndpoint("https://updates.push.services.mozilla.com/wpush/v2/abc"))
	assert.NoError(t, ValidatePushEndpoint("https://fcm.googleapis.com/fcm/send/xyz"))

	for _, bad := range []string{
		"",
		"http://fcm.googleapis.com/fcm/send/xyz",
		"file:///etc/passwd",
		"https://localhost/push",
		"https://sub.localhost/push",
		"https://printer.local/push",
		"https://metadata.internal/push",
		"https://127.0.0.1/push",
		"https://169.254.169.254/latest/meta-data/",
		"https://10.0.0.5/push",
		"https://192.168.1.1/push",
		"https://100.64.0.1/push",
		"https://[::1]/push",
		"https://user:pass@fcm.googleapis.com/push",
	} {
		assert.Error(t, ValidatePushEndpoint(bad), "expected error for %q", bad)
	}
}

func TestSafeWebPushHTTPClientRejectsLoopback(t *testing.T) {
	client := newSafeWebPushHTTPClient(time.Second)
	_, err := client.Get("https://127.0.0.1:1/push")
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "non-public IP"), "unexpected error: %v", err)
}

func TestWebPushServiceNotifyHeadersAndErrorBody(t *testing.T) {
	priv, pub, p256dh, auth := testKeys(t)

	var gotTTL string
	var gotAuth string
	statusToReturn := http.StatusCreated
	bodyToReturn := ""

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTTL = r.Header.Get("TTL")
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(statusToReturn)
		if bodyToReturn != "" {
			_, _ = w.Write([]byte(bodyToReturn))
		}
	}))
	defer srv.Close()

	// Point subscriber config through NewWebPushService to verify normalization.
	subAddr := "mailto:admin@example.com"
	svc, err := NewWebPushService(e2clickermodule.Notification{
		ClientTimeout: "5s",
		WebPush: &e2clickermodule.WebPushJSON{
			Value: e2clickermodule.WebPushSubmodule{
				PrivateKey: priv,
				PublicKey:  pub,
				Subscriber: &subAddr,
			},
		},
	})
	assert.NoError(t, err)
	svc.http = srv.Client()

	// Call send directly via webpush-go with httptest TLS endpoint (bypassing public-IP
	// endpoint check only to test wire headers & body reading against local httptest.Server).
	// Verify Notify rejects loopback endpoints first:
	sub := openapi.PushSubscription{
		DeviceID: "dev1",
		Endpoint: srv.URL,
	}
	sub.Keys.P256Dh = p256dh
	sub.Keys.Auth = auth

	err = svc.Notify(context.Background(), Notification{
		Type:     openapi.TestMessage,
		Username: "Happy Sunflower",
		Message:  openapi.NotificationMessage{Title: "Hi", Message: "Test"},
	}, sub)
	assert.Error(t, err)

	// Now verify actual VAPID JWT & TTL & resp.Body read using webpush.SendNotificationWithContext
	// with the exact same options constructed by Notify:
	resp, err := webpush.SendNotificationWithContext(
		context.Background(),
		[]byte(`{"type":"test_message","username":"Happy Sunflower","message":{"title":"Hi","message":"Test"}}`),
		convertSubscription(sub),
		&webpush.Options{
			HTTPClient:      svc.http,
			TTL:             webPushTTLSeconds,
			Urgency:         webpush.UrgencyHigh,
			Subscriber:      svc.subscriber,
			VAPIDPublicKey:  svc.config.PublicKey,
			VAPIDPrivateKey: svc.config.PrivateKey,
		},
	)
	assert.NoError(t, err)
	resp.Body.Close()

	assert.Equal(t, "86400", gotTTL)
	assert.True(t, strings.HasPrefix(gotAuth, "vapid t="))

	jwtPart := strings.TrimPrefix(strings.Split(gotAuth, ",")[0], "vapid t=")
	parts := strings.Split(jwtPart, ".")
	assert.Equal(t, 3, len(parts))
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	assert.NoError(t, err)

	var claims map[string]any
	assert.NoError(t, json.Unmarshal(payloadBytes, &claims))
	assert.Equal(t, "mailto:admin@example.com", claims["sub"])

	// Verify consumeHTTPUnknownStatusError reads body and is NOT marked public.
	statusToReturn = http.StatusGone
	bodyToReturn = "push subscription expired"
	resp, err = webpush.SendNotificationWithContext(
		context.Background(),
		[]byte(`{}`),
		convertSubscription(sub),
		&webpush.Options{
			HTTPClient:      svc.http,
			TTL:             webPushTTLSeconds,
			Urgency:         webpush.UrgencyHigh,
			Subscriber:      svc.subscriber,
			VAPIDPublicKey:  svc.config.PublicKey,
			VAPIDPrivateKey: svc.config.PrivateKey,
		},
	)
	assert.NoError(t, err)
	defer resp.Body.Close()

	httpErr := consumeHTTPUnknownStatusError(resp)
	assert.True(t, strings.Contains(httpErr.Error(), "push subscription expired"))
	assert.False(t, publicerrors.IsPublic(httpErr), "HTTPUnknownStatusError must not be public")
}

type memUserNotificationStorage struct {
	mu    sync.Mutex
	prefs map[user.Secret]UserPreferences
}

func (m *memUserNotificationStorage) UserPreferences(ctx context.Context, secret user.Secret) (UserPreferences, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.prefs[secret], nil
}

func (m *memUserNotificationStorage) SetUserPreferencesTx(ctx context.Context, secret user.Secret, set func(*UserPreferences) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	p := m.prefs[secret]
	if err := set(&p); err != nil {
		return err
	}
	m.prefs[secret] = p
	return nil
}

func TestUserNotificationServicePreservationAndCaps(t *testing.T) {
	_, _, p256dh, auth := testKeys(t)
	store := &memUserNotificationStorage{prefs: map[user.Secret]UserPreferences{}}

	svc := &userNotificationService{
		userNotifications: store,
		logger:            slogt.New(t),
	}

	ctx := context.Background()
	secret := user.Secret("secret1")

	// Empty channels must return ErrNoNotificationChannels.
	err := svc.NotifyUser(ctx, secret, openapi.TestMessage)
	assert.True(t, errors.Is(err, ErrNoNotificationChannels))
	assert.True(t, publicerrors.IsPublic(err))

	// Subscribe device 1.
	sub := func(id string) openapi.PushSubscription {
		s := openapi.PushSubscription{
			DeviceID: id,
			Endpoint: "https://fcm.googleapis.com/fcm/send/" + id,
		}
		s.Keys.P256Dh = p256dh
		s.Keys.Auth = auth
		return s
	}
	assert.NoError(t, svc.SubscribeWebPush(ctx, secret, sub("dev1")))

	// Set Email preferences; WebPush subscription must be preserved.
	assert.NoError(t, svc.SetUserPreferences(ctx, secret, &UserPreferences{
		NotificationConfigs: NotificationConfigs{
			Email: []EmailNotificationConfig{{Address: "user@example.com"}},
		},
	}))

	prefs, err := svc.UserPreferences(ctx, secret)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(prefs.NotificationConfigs.Email))
	assert.Equal(t, 1, len(prefs.NotificationConfigs.WebPush))
	assert.Equal(t, "dev1", prefs.NotificationConfigs.WebPush[0].DeviceID)

	// Adding 15 more succeeds (total 16), 17th fails with ErrTooManyWebPushSubscriptions.
	for i := 2; i <= maxWebPushSubscriptionsPerUser; i++ {
		assert.NoError(t, svc.SubscribeWebPush(ctx, secret, sub("dev"+string(rune('A'+i)))))
	}
	err = svc.SubscribeWebPush(ctx, secret, sub("dev17"))
	assert.True(t, errors.Is(err, ErrTooManyWebPushSubscriptions))

	// Unsubscribing dev1 removes it while keeping email.
	assert.NoError(t, svc.UnsubscribeWebPush(ctx, secret, "dev1"))
	prefs, err = svc.UserPreferences(ctx, secret)
	assert.NoError(t, err)
	assert.Equal(t, maxWebPushSubscriptionsPerUser-1, len(prefs.NotificationConfigs.WebPush))
	assert.Equal(t, 1, len(prefs.NotificationConfigs.Email))
}
