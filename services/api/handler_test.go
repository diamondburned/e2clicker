package api

import (
	"context"
	"encoding/json"
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

type stubUserStorage struct{}

func (stubUserStorage) CreateUser(context.Context, user.Secret, string) (user.User, error) {
	return user.User{}, nil
}

func (stubUserStorage) User(context.Context, user.Secret) (user.User, error) {
	return user.User{}, user.ErrUnknownUser
}

func (stubUserStorage) UpdateUserName(context.Context, user.Secret, string) error {
	return nil
}

func (stubUserStorage) UpdateUserLocale(context.Context, user.Secret, user.Locale) error {
	return nil
}

type stubUserSessionStorage struct {
	registerErr error
}

func (s stubUserSessionStorage) RegisterSession(context.Context, []byte, user.Secret, string) error {
	return s.registerErr
}

func (stubUserSessionStorage) ValidateSession(context.Context, []byte) (user.Session, error) {
	return user.Session{}, user.ErrInvalidSession
}

func (stubUserSessionStorage) ListSessions(context.Context, user.Secret) ([]user.Session, error) {
	return nil, nil
}

func (stubUserSessionStorage) DeleteSession(context.Context, user.Secret, int64) error {
	return nil
}

func TestAuthUnknownUser(t *testing.T) {
	users, err := user.NewUserService(user.UserServiceConfig{
		UserStorage:        stubUserStorage{},
		UserSessionStorage: stubUserSessionStorage{registerErr: user.ErrUnknownUser},
	})
	assert.NoError(t, err)

	lc := fxtest.NewLifecycle(t)
	h := newOpenAPIHandler(OpenAPIHandlerServices{
		Users: users,
	}, slogt.New(t), lc)

	ctx := context.Background()
	_, err = h.Auth(ctx, openapi.AuthRequestObject{
		Body: &openapi.AuthJSONRequestBody{Secret: "NONEXISTENTUSER1"},
	})
	assert.Error(t, err)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth", nil).WithContext(ctx)
	writeError(rec, req, err, 0)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var body openapi.Error
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "unknown user", body.Message)
	assert.Equal(t, (*bool)(nil), body.Internal)
}
