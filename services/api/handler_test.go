package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"e2clicker.app/internal/ptr"
	"e2clicker.app/services/api/openapi"
	"e2clicker.app/services/notification"
	notificationapi "e2clicker.app/services/notification/openapi"
	"e2clicker.app/services/user"
	"github.com/alecthomas/assert/v2"
	"github.com/neilotoole/slogt"
	"go.uber.org/fx/fxtest"
	"libdb.so/ctxt"
)

type stubUserNotificationService struct {
	notifyErr      error
	lastAttempt    *notificationapi.NotificationAttempt
	lastAttemptErr error
}

func (s *stubUserNotificationService) NotifyUser(ctx context.Context, secret user.Secret, t notificationapi.NotificationType) error {
	return s.notifyErr
}

func (s *stubUserNotificationService) UserPreferences(ctx context.Context, secret user.Secret) (notification.UserPreferences, error) {
	return notification.UserPreferences{}, nil
}

func (s *stubUserNotificationService) LastNotificationAttempt(ctx context.Context, secret user.Secret) (*notificationapi.NotificationAttempt, error) {
	return s.lastAttempt, s.lastAttemptErr
}

func (s *stubUserNotificationService) SetUserPreferences(ctx context.Context, secret user.Secret, preferences *notification.UserPreferences) error {
	return nil
}

func (s *stubUserNotificationService) SetUserPreferencesSafe(ctx context.Context, secret user.Secret, newPreferences, oldPreferences *notification.UserPreferences) error {
	return nil
}

func (s *stubUserNotificationService) SubscribeWebPush(ctx context.Context, secret user.Secret, subscription notificationapi.PushSubscription) error {
	return nil
}

func (s *stubUserNotificationService) UnsubscribeWebPush(ctx context.Context, secret user.Secret, deviceID string) error {
	return nil
}

func TestSendTestNotificationRateLimitAndErrors(t *testing.T) {
	notifs := &stubUserNotificationService{}
	lc := fxtest.NewLifecycle(t)
	h := newOpenAPIHandler(OpenAPIHandlerServices{
		UserNotifications: notifs,
	}, slogt.New(t), lc)

	ctx := ctxt.With(context.Background(), user.Session{UserSecret: "user-secret-1"})

	// ErrNoNotificationChannels must return public 400.
	notifs.notifyErr = notification.ErrNoNotificationChannels
	_, err := h.SendTestNotification(ctx, openapi.SendTestNotificationRequestObject{})
	assert.Error(t, err)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/notifications/test", nil).WithContext(ctx)
	writeError(rec, req, err, 0)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Upstream HTTPUnknownStatusError must not leak response body and must return 500.
	notifs.notifyErr = notification.HTTPUnknownStatusError{
		StatusCode: 502,
		Body:       "internal-secret-backend-body",
	}
	_, err = h.SendTestNotification(ctx, openapi.SendTestNotificationRequestObject{})
	assert.Error(t, err)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/notifications/test", nil).WithContext(ctx)
	writeError(rec, req, err, 0)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.False(t, strings.Contains(rec.Body.String(), "internal-secret-backend-body"))

	// Burst of 3 succeeds (2 tokens spent above -> 3rd succeeds, 4th fails with 429).
	notifs.notifyErr = nil
	_, err = h.SendTestNotification(ctx, openapi.SendTestNotificationRequestObject{})
	assert.NoError(t, err)

	_, err = h.SendTestNotification(ctx, openapi.SendTestNotificationRequestObject{})
	assert.Error(t, err)

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/notifications/test", nil).WithContext(ctx)
	writeError(rec, req, err, 0)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.NotEqual(t, "", rec.Header().Get("Retry-After"))
}

func TestUserNotificationStatus(t *testing.T) {
	notifs := &stubUserNotificationService{}
	lc := fxtest.NewLifecycle(t)
	h := newOpenAPIHandler(OpenAPIHandlerServices{
		UserNotifications: notifs,
	}, slogt.New(t), lc)

	ctx := ctxt.With(context.Background(), user.Session{UserSecret: "user-secret-1"})

	resp, err := h.UserNotificationStatus(ctx, openapi.UserNotificationStatusRequestObject{})
	assert.NoError(t, err)
	statusResp, ok := resp.(openapi.UserNotificationStatus200JSONResponse)
	assert.True(t, ok)
	assert.Equal(t, (*openapi.NotificationAttempt)(nil), statusResp.LastAttempt)

	sentAt := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	notifs.lastAttempt = &notificationapi.NotificationAttempt{
		SentAt:      sentAt,
		Errored:     true,
		ErrorReason: ptr.To("webPush: subscription expired"),
	}

	resp, err = h.UserNotificationStatus(ctx, openapi.UserNotificationStatusRequestObject{})
	assert.NoError(t, err)
	statusResp, ok = resp.(openapi.UserNotificationStatus200JSONResponse)
	assert.True(t, ok)
	assert.True(t, statusResp.LastAttempt != nil)
	assert.Equal(t, sentAt, statusResp.LastAttempt.SentAt)
	assert.True(t, statusResp.LastAttempt.Errored)
	assert.Equal(t, ptr.To("webPush: subscription expired"), statusResp.LastAttempt.ErrorReason)
}
