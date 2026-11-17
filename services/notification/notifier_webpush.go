package notification

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"strings"
	"time"

	"e2clicker.app/services/notification/openapi"
	"github.com/SherClockHolmes/webpush-go"

	e2clickermodule "e2clicker.app/nix/modules/e2clicker"
)

const (
	maxWebPushSubscriptionsPerUser = 16
	webPushTTLSeconds              = 24 * 60 * 60 // 24h
)

var cgnatCIDR = func() *net.IPNet {
	_, ipNet, err := net.ParseCIDR("100.64.0.0/10")
	if err != nil {
		panic(err)
	}
	return ipNet
}()

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
	http       *http.Client
	config     *e2clickermodule.WebPushSubmodule
	subscriber string
}

// NewWebPushService creates a new Push API service.
func NewWebPushService(config e2clickermodule.Notification) (*WebPushService, error) {
	if config.WebPush == nil {
		return nil, nil
	}

	timeout, err := time.ParseDuration(config.ClientTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid client timeout %q: %w", config.ClientTimeout, err)
	}

	var keys *e2clickermodule.WebPushSubmodule

	switch value := config.WebPush.Value.(type) {
	case e2clickermodule.WebPushSubmodule:
		keys = &value
	case e2clickermodule.WebPushPath:
		b, err := os.ReadFile(string(value))
		if err != nil {
			return nil, fmt.Errorf("cannot read WebPushKeys file at %s: %w", value, err)
		}
		keys = new(e2clickermodule.WebPushSubmodule)
		if err := json.Unmarshal(b, keys); err != nil {
			return nil, fmt.Errorf("cannot unmarshal WebPushKeys at %s: %w", value, err)
		}
	default:
		panic("unreachable")
	}

	var rawSubscriber string
	if keys.Subscriber != nil && *keys.Subscriber != "" {
		rawSubscriber = *keys.Subscriber
	} else if config.WebPushSubscriber != nil {
		rawSubscriber = *config.WebPushSubscriber
	}

	subscriber, err := normalizeVAPIDSubscriber(rawSubscriber)
	if err != nil {
		return nil, err
	}

	return &WebPushService{
		http:       newSafeWebPushHTTPClient(timeout),
		config:     keys,
		subscriber: subscriber,
	}, nil
}

// VAPIDPublicKey returns the VAPID public key.
func (s WebPushService) VAPIDPublicKey() string {
	return s.config.PublicKey
}

func (s WebPushService) Notify(ctx context.Context, n Notification, config WebPushNotificationConfig) error {
	if !config.ExpirationTime.IsZero() && config.ExpirationTime.Before(time.Now()) {
		return &WebPushSubscriptionExpired{config.ExpirationTime}
	}

	if err := ValidatePushSubscription(config); err != nil {
		return ConfigError{Service: "webPush", err: err}
	}

	m, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("cannot marshal notification: %w", err)
	}

	opts := &webpush.Options{
		HTTPClient:      s.http,
		TTL:             webPushTTLSeconds,
		Urgency:         webpush.UrgencyHigh,
		Subscriber:      s.subscriber,
		VAPIDPublicKey:  s.config.PublicKey,
		VAPIDPrivateKey: s.config.PrivateKey,
	}

	resp, err := webpush.SendNotificationWithContext(ctx, m, convertSubscription(config), opts)
	if err != nil {
		return fmt.Errorf("cannot send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return consumeHTTPUnknownStatusError(resp)
	}

	return nil
}

// ValidatePushSubscription validates a Web Push subscription payload.
func ValidatePushSubscription(sub openapi.PushSubscription) error {
	if sub.DeviceID == "" || len(sub.DeviceID) > 64 || strings.ContainsAny(sub.DeviceID, " \t\r\n/") {
		return errors.New("invalid deviceID")
	}
	if err := ValidatePushEndpoint(sub.Endpoint); err != nil {
		return err
	}
	if _, err := decodeBase64Key(sub.Keys.P256Dh, 16); err != nil {
		return fmt.Errorf("invalid keys.p256dh: %w", err)
	}
	if _, err := decodeBase64Key(sub.Keys.Auth, 16); err != nil {
		return fmt.Errorf("invalid keys.auth: %w", err)
	}
	return nil
}

// ValidatePushEndpoint checks that a Web Push endpoint is a valid public HTTPS URL.
func ValidatePushEndpoint(rawURL string) error {
	if rawURL == "" || len(rawURL) > 2048 {
		return errors.New("endpoint must be between 1 and 2048 characters")
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid endpoint URL: %w", err)
	}
	if u.Scheme != "https" {
		return errors.New("endpoint URL must use https://")
	}
	if u.User != nil {
		return errors.New("endpoint URL must not include credentials")
	}

	host := strings.ToLower(u.Hostname())
	if host == "" {
		return errors.New("endpoint URL must include a host")
	}
	if host == "localhost" ||
		strings.HasSuffix(host, ".localhost") ||
		strings.HasSuffix(host, ".local") ||
		strings.HasSuffix(host, ".internal") {
		return errors.New("endpoint host is not allowed")
	}
	if ip := net.ParseIP(host); ip != nil && !isPublicIP(ip) {
		return errors.New("endpoint IP address is not allowed")
	}

	return nil
}

func decodeBase64Key(raw string, minBytes int) ([]byte, error) {
	if raw == "" || len(raw) > 256 {
		return nil, errors.New("key length out of bounds")
	}

	encodings := []*base64.Encoding{
		base64.RawURLEncoding,
		base64.URLEncoding,
		base64.RawStdEncoding,
		base64.StdEncoding,
	}
	for _, enc := range encodings {
		if b, err := enc.DecodeString(raw); err == nil {
			if len(b) < minBytes {
				return nil, fmt.Errorf("decoded key too short (%d bytes)", len(b))
			}
			return b, nil
		}
	}
	return nil, errors.New("invalid base64 encoding")
}

func normalizeVAPIDSubscriber(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("webPush subscriber must be set (e.g. mailto:... or https://...)")
	}
	if strings.HasPrefix(raw, "https://") {
		u, err := url.ParseRequestURI(raw)
		if err != nil || u.Host == "" || u.User != nil {
			return "", fmt.Errorf("invalid https webPush subscriber %q", raw)
		}
		return raw, nil
	}
	if strings.HasPrefix(raw, "http://") {
		return "", fmt.Errorf("webPush subscriber must use mailto: or https://, got %q", raw)
	}
	addr := strings.TrimPrefix(raw, "mailto:")
	parsed, err := mail.ParseAddress(addr)
	if err != nil || parsed.Address == "" || !strings.Contains(parsed.Address, "@") {
		return "", fmt.Errorf("invalid mailto webPush subscriber %q", raw)
	}
	// webpush-go automatically prepends "mailto:" to non-https subscribers.
	return parsed.Address, nil
}

func isPublicIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		ip = ip4
	}
	if !ip.IsGlobalUnicast() ||
		ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() ||
		ip.IsUnspecified() {
		return false
	}
	if cgnatCIDR.Contains(ip) {
		return false
	}
	return true
}

func newSafeWebPushHTTPClient(timeout time.Duration) *http.Client {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		ForceAttemptHTTP2: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, err
			}
			for _, ipAddr := range ips {
				if !isPublicIP(ipAddr.IP) {
					return nil, fmt.Errorf("refusing to dial non-public IP %s for host %s", ipAddr.IP, host)
				}
			}
			var lastErr error
			for _, ipAddr := range ips {
				conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(ipAddr.IP.String(), port))
				if err == nil {
					return conn, nil
				}
				lastErr = err
			}
			if lastErr != nil {
				return nil, lastErr
			}
			return nil, fmt.Errorf("no IP addresses resolved for %s", host)
		},
	}
	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("redirects are not allowed for web push endpoints")
		},
	}
}
