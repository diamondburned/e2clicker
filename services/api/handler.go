package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"time"
	"uuid"

	"e2clicker.app/internal/ptr"
	"e2clicker.app/internal/publicerrors"
	"e2clicker.app/internal/userlimit"
	"e2clicker.app/services/api/openapi"
	"e2clicker.app/services/dosage"
	"e2clicker.app/services/notification"
	"e2clicker.app/services/user"
	"go.uber.org/fx"
	"golang.org/x/time/rate"

	dosageapi "e2clicker.app/services/dosage/openapi"
	notificationapi "e2clicker.app/services/notification/openapi"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// openAPIHandler is the handler for the OpenAPI service.
// It implements the OpenAPI service interface.
type openAPIHandler struct {
	logger           *slog.Logger
	users            *user.UserService
	notifs           notification.UserNotificationService
	notif            notification.NotificationService
	dosage           dosage.DosageStorage
	doseHistory      dosage.DoseHistoryStorage
	testNotifLimiter *userlimit.UserRateLimiter[user.Secret]
}

// OpenAPIHandlerServices is the set of service dependencies required by the
// OpenAPIHandler.
type OpenAPIHandlerServices struct {
	fx.In

	Users             *user.UserService
	Notification      notification.NotificationService
	UserNotifications notification.UserNotificationService
	Dosage            dosage.DosageStorage
	DoseHistory       dosage.DoseHistoryStorage
}

// newOpenAPIHandler creates a new OpenAPIHandler.
func newOpenAPIHandler(deps OpenAPIHandlerServices, logger *slog.Logger, lc fx.Lifecycle) *openAPIHandler {
	limiter := userlimit.NewUserRateLimiter[user.Secret](rate.Every(1*time.Minute), 3)
	stopCleanup := limiter.BeginCleanup()
	lc.Append(fx.StopHook(func(ctx context.Context) error {
		stopCleanup()
		return nil
	}))

	return &openAPIHandler{
		logger:           logger,
		users:            deps.Users,
		notifs:           deps.UserNotifications,
		notif:            deps.Notification,
		dosage:           deps.Dosage,
		doseHistory:      deps.DoseHistory,
		testNotifLimiter: limiter,
	}
}

func (h *openAPIHandler) asHandler() openapi.ServerInterface {
	return openapi.NewStrictHandlerWithOptions(
		h, nil,
		openapi.StrictHTTPServerOptions{
			RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				err = publicerrors.ForcePublic(err) // only validation errors
				writeError(w, r, err, http.StatusBadRequest)
			},
			ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				writeError(w, r, err, 0)
			},
		},
	)
}

// Register a new account
// (POST /register)
func (h *openAPIHandler) Register(ctx context.Context, request openapi.RegisterRequestObject) (openapi.RegisterResponseObject, error) {
	u, err := h.users.CreateUser(ctx, request.Body.Name)
	if err != nil {
		return nil, err
	}

	return openapi.Register200JSONResponse{
		Name:   u.Name,
		Locale: u.Locale,
		Secret: u.Secret,
	}, nil
}

// Authenticate a user
// (POST /auth)
func (h *openAPIHandler) Auth(ctx context.Context, request openapi.AuthRequestObject) (openapi.AuthResponseObject, error) {
	t, err := h.users.CreateSession(ctx, request.Body.Secret, optstr(request.Params.UserAgent))
	if err != nil {
		return nil, err
	}

	return openapi.Auth200JSONResponse{
		Token: string(t),
	}, nil
}

// Get the current user
// (GET /me)
func (h *openAPIHandler) CurrentUser(ctx context.Context, request openapi.CurrentUserRequestObject) (openapi.CurrentUserResponseObject, error) {
	session := sessionFromCtx(ctx)

	u, err := h.users.User(ctx, session.UserSecret)
	if err != nil {
		return nil, err
	}

	return openapi.CurrentUser200JSONResponse{
		Name:   u.Name,
		Locale: u.Locale,
		Secret: session.UserSecret,
	}, nil
}

// List the current user's sessions
// (GET /me/sessions)
func (h *openAPIHandler) CurrentUserSessions(ctx context.Context, request openapi.CurrentUserSessionsRequestObject) (openapi.CurrentUserSessionsResponseObject, error) {
	session := sessionFromCtx(ctx)

	s, err := h.users.ListSessions(ctx, session.UserSecret)
	if err != nil {
		return nil, err
	}

	return openapi.CurrentUserSessions200JSONResponse(convertList(s, convertSession)), nil
}

// Delete one of the current user's sessions
// (DELETE /me/sessions)
func (h *openAPIHandler) DeleteUserSession(ctx context.Context, request openapi.DeleteUserSessionRequestObject) (openapi.DeleteUserSessionResponseObject, error) {
	session := sessionFromCtx(ctx)

	err := h.users.DeleteSession(ctx, session.UserSecret, request.Params.ID)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteUserSession204Response{}, nil
}

// List all available delivery methods
// (GET /delivery-methods)
func (h *openAPIHandler) DeliveryMethods(ctx context.Context, request openapi.DeliveryMethodsRequestObject) (openapi.DeliveryMethodsResponseObject, error) {
	methods, err := h.dosage.DeliveryMethods(ctx)
	if err != nil {
		return nil, err
	}

	return openapi.DeliveryMethods200JSONResponse(
		convertList(methods, func(m dosage.DeliveryMethod) openapi.DeliveryMethod {
			return openapi.DeliveryMethod(m)
		}),
	), nil
}

func (h *openAPIHandler) SetDosage(ctx context.Context, request openapi.SetDosageRequestObject) (openapi.SetDosageResponseObject, error) {
	session := sessionFromCtx(ctx)

	methods, err := h.dosage.DeliveryMethods(ctx)
	if err != nil {
		return nil, err
	}

	s := *(*dosageapi.Dosage)(request.Body)
	if !slices.ContainsFunc(methods, func(m dosage.DeliveryMethod) bool {
		return m.ID == s.DeliveryMethod
	}) {
		return nil, publicerrors.Errorf("invalid delivery method %q", s.DeliveryMethod)
	}

	if err := h.dosage.SetDosage(ctx, session.UserSecret, s); err != nil {
		return nil, err
	}

	return openapi.SetDosage204Response{}, nil
}

func (h *openAPIHandler) ClearDosage(ctx context.Context, request openapi.ClearDosageRequestObject) (openapi.ClearDosageResponseObject, error) {
	session := sessionFromCtx(ctx)
	if err := h.dosage.ClearDosage(ctx, session.UserSecret); err != nil {
		return nil, err
	}
	return openapi.ClearDosage204Response{}, nil
}

func (h *openAPIHandler) RecordDose(ctx context.Context, request openapi.RecordDoseRequestObject) (openapi.RecordDoseResponseObject, error) {
	session := sessionFromCtx(ctx)

	now := time.Now().Truncate(time.Second)

	d, err := h.dosage.Dosage(ctx, session.UserSecret)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, publicerrors.New("no dosage set")
	}

	dose := dosage.Dose{
		DeliveryMethod: d.DeliveryMethod,
		Dose:           d.Dose,
		TakenAt:        ptr.DerefOr(request.Body.TakenAt, now),
		Comment:        ptr.Deref(request.Body.Comment),
	}

	recorded, err := h.doseHistory.RecordDose(ctx, session.UserSecret, dose)
	if err != nil {
		return nil, err
	}

	return openapi.RecordDose200JSONResponse(openapi.Dose(recorded.ToOpenAPI())), nil
}

func (h *openAPIHandler) EditDose(ctx context.Context, request openapi.EditDoseRequestObject) (openapi.EditDoseResponseObject, error) {
	session := sessionFromCtx(ctx)

	o := dosage.Dose{
		DeliveryMethod: request.Body.DeliveryMethod,
		Dose:           request.Body.Dose,
		TakenAt:        request.Body.TakenAt,
		TakenOffAt:     request.Body.TakenOffAt,
		Comment:        ptr.Deref(request.Body.Comment),
	}

	if err := h.doseHistory.EditDose(ctx, session.UserSecret, uuid.UUID(request.DoseID), o); err != nil {
		return nil, err
	}

	return openapi.EditDose204Response{}, nil
}

func (h *openAPIHandler) ForgetDose(ctx context.Context, request openapi.ForgetDoseRequestObject) (openapi.ForgetDoseResponseObject, error) {
	session := sessionFromCtx(ctx)
	if err := h.doseHistory.ForgetDoses(ctx, session.UserSecret, []uuid.UUID{uuid.UUID(request.DoseID)}); err != nil {
		return nil, err
	}
	return openapi.ForgetDose204Response{}, nil
}

func (h *openAPIHandler) ForgetDoses(ctx context.Context, request openapi.ForgetDosesRequestObject) (openapi.ForgetDosesResponseObject, error) {
	session := sessionFromCtx(ctx)
	ids := convertList(request.Params.Ids, func(id openapi_types.UUID) uuid.UUID {
		return uuid.UUID(id)
	})
	if err := h.doseHistory.ForgetDoses(ctx, session.UserSecret, ids); err != nil {
		return nil, err
	}
	return openapi.ForgetDoses204Response{}, nil
}

func (h *openAPIHandler) Dosage(ctx context.Context, request openapi.DosageRequestObject) (openapi.DosageResponseObject, error) {
	session := sessionFromCtx(ctx)

	var r openapi.Dosage200JSONResponse

	dosage, err := h.dosage.Dosage(ctx, session.UserSecret)
	if err != nil {
		return nil, fmt.Errorf("cannot get dosage: %w", err)
	}
	if dosage != nil {
		r.Dosage = (*openapi.Dosage)(dosage)
	}

	if request.Params.Start != nil {
		end := ptr.DerefOr(request.Params.End, time.Now())

		const oneYear = 365 * 24 * time.Hour
		if end.Sub(*request.Params.Start) > oneYear {
			return nil, publicerrors.New("" +
				"requested history range is too large, must be 1 year or less " +
				"(consider exporting instead)")
		}

		os := make([]openapi.Dose, 0, 32)
		r.History = &os

		doses, err := h.doseHistory.DoseHistory(
			ctx, session.UserSecret,
			*request.Params.Start,
			end)
		if err != nil {
			return nil, fmt.Errorf("cannot get dosage history: %w", err)
		}

		for _, dose := range doses {
			os = append(os, openapi.Dose(dose.ToOpenAPI()))
		}
	}

	return r, nil
}

func (h openAPIHandler) ExportDoses(ctx context.Context, request openapi.ExportDosesRequestObject) (openapi.ExportDosesResponseObject, error) {
	panic("unreachable") // see handler_importexport.go
}

func (h *openAPIHandler) ImportDoses(ctx context.Context, request openapi.ImportDosesRequestObject) (openapi.ImportDosesResponseObject, error) {
	panic("unreachable") // see handler_importexport.go
}

func (h *openAPIHandler) WebPushInfo(ctx context.Context, request openapi.WebPushInfoRequestObject) (openapi.WebPushInfoResponseObject, error) {
	config := h.notif.Config()
	if config.WebPush == nil {
		return nil, notification.ErrWebPushNotAvailable
	}
	return openapi.WebPushInfo200JSONResponse(openapi.PushInfo{
		ApplicationServerKey: config.WebPush.VAPIDPublicKey(),
	}), nil
}

// Get the server's supported notification methods
// (GET /notifications/methods)
func (h *openAPIHandler) SupportedNotificationMethods(ctx context.Context, request openapi.SupportedNotificationMethodsRequestObject) (openapi.SupportedNotificationMethodsResponseObject, error) {
	addIfTrue := func(ret []string, b bool, value string) []string {
		if b {
			ret = append(ret, value)
		}
		return ret
	}

	supports := h.notif.Config().Supports()

	var ret openapi.NotificationMethodSupports
	ret = addIfTrue(ret, supports.Gotify, "gotify")
	ret = addIfTrue(ret, supports.Pushover, "pushover")
	ret = addIfTrue(ret, supports.WebPush, "webPush")
	ret = addIfTrue(ret, supports.Email, "email")

	return openapi.SupportedNotificationMethods200JSONResponse(openapi.NotificationMethodSupports(ret)), nil
}

// Get the user's notification preferences
// (GET /notifications/preferences)
func (h *openAPIHandler) UserNotificationPreferences(ctx context.Context, request openapi.UserNotificationPreferencesRequestObject) (openapi.UserNotificationPreferencesResponseObject, error) {
	session := sessionFromCtx(ctx)

	p, err := h.notifs.UserPreferences(ctx, session.UserSecret)
	if err != nil {
		return nil, err
	}

	var ret openapi.NotificationPreferences

	if len(p.CustomNotifications) > 0 {
		ret.CustomNotifications = make(map[string]openapi.NotificationMessage, len(p.CustomNotifications))
		for k, v := range p.CustomNotifications {
			ret.CustomNotifications[k] = openapi.NotificationMessage(v)
		}
	}

	if len(p.NotificationConfigs.WebPush) > 0 {
		s := make([]openapi.PushDevice, len(p.NotificationConfigs.WebPush))
		for i, sub := range p.NotificationConfigs.WebPush {
			s[i] = openapi.PushDevice{
				DeviceID:       sub.DeviceID,
				ExpirationTime: sub.ExpirationTime,
			}
		}
		ret.NotificationConfigs.WebPush = &s
	}

	if len(p.NotificationConfigs.Email) > 0 {
		s := make([]openapi.EmailSubscription, len(p.NotificationConfigs.Email))
		for i, sub := range p.NotificationConfigs.Email {
			s[i] = openapi.EmailSubscription(sub)
		}
		ret.NotificationConfigs.Email = &s
	}

	return openapi.UserNotificationPreferences200JSONResponse(ret), nil
}

// Get the user's notification status
// (GET /notifications/status)
func (h *openAPIHandler) UserNotificationStatus(ctx context.Context, request openapi.UserNotificationStatusRequestObject) (openapi.UserNotificationStatusResponseObject, error) {
	session := sessionFromCtx(ctx)

	attempt, err := h.notifs.LastNotificationAttempt(ctx, session.UserSecret)
	if err != nil {
		return nil, err
	}

	var ret openapi.NotificationStatus
	if attempt != nil {
		ret.LastAttempt = (*openapi.NotificationAttempt)(attempt)
	}

	return openapi.UserNotificationStatus200JSONResponse(ret), nil
}

// Update the user's notification preferences
// (PUT /notifications/preferences)
func (h *openAPIHandler) UserUpdateNotificationPreferences(ctx context.Context, request openapi.UserUpdateNotificationPreferencesRequestObject) (openapi.UserUpdateNotificationPreferencesResponseObject, error) {
	session := sessionFromCtx(ctx)

	newPreferences := &notification.UserPreferences{}

	if request.Body.CustomNotifications != nil {
		newPreferences.CustomNotifications = make(notificationapi.CustomNotifications, len(request.Body.CustomNotifications))
		for k, v := range request.Body.CustomNotifications {
			newPreferences.CustomNotifications[k] = notificationapi.NotificationMessage(v)
		}
	}

	if request.Body.NotificationConfigs.Email != nil {
		newPreferences.NotificationConfigs.Email = make([]notificationapi.EmailSubscription, len(*request.Body.NotificationConfigs.Email))
		for i, v := range *request.Body.NotificationConfigs.Email {
			newPreferences.NotificationConfigs.Email[i] = notificationapi.EmailSubscription(v)
		}
	}

	if err := h.notifs.SetUserPreferences(ctx, session.UserSecret, newPreferences); err != nil {
		return nil, err
	}

	return openapi.UserUpdateNotificationPreferences204Response{}, nil
}

// Subscribe or update a Web Push device subscription
// (PUT /notifications/web-push/{deviceID})
func (h *openAPIHandler) UserSubscribePush(ctx context.Context, request openapi.UserSubscribePushRequestObject) (openapi.UserSubscribePushResponseObject, error) {
	if h.notif.Config().WebPush == nil {
		return nil, notification.ErrWebPushNotAvailable
	}
	if request.Body.DeviceID != request.DeviceID {
		return nil, publicerrors.New("deviceID in path does not match request body")
	}

	session := sessionFromCtx(ctx)
	if err := h.notifs.SubscribeWebPush(ctx, session.UserSecret, notificationapi.PushSubscription(*request.Body)); err != nil {
		return nil, err
	}

	return openapi.UserSubscribePush204Response{}, nil
}

// Unsubscribe a Web Push device subscription
// (DELETE /notifications/web-push/{deviceID})
func (h *openAPIHandler) UserUnsubscribePush(ctx context.Context, request openapi.UserUnsubscribePushRequestObject) (openapi.UserUnsubscribePushResponseObject, error) {
	session := sessionFromCtx(ctx)
	if err := h.notifs.UnsubscribeWebPush(ctx, session.UserSecret, request.DeviceID); err != nil {
		return nil, err
	}

	return openapi.UserUnsubscribePush204Response{}, nil
}

// Send a test notification
// (POST /notifications/test)
func (h *openAPIHandler) SendTestNotification(ctx context.Context, request openapi.SendTestNotificationRequestObject) (openapi.SendTestNotificationResponseObject, error) {
	session := sessionFromCtx(ctx)

	limit := h.testNotifLimiter.Reserve(session.UserSecret)
	if err := userlimit.AsError(limit); err != nil {
		return nil, err
	}

	if err := h.notifs.NotifyUser(ctx, session.UserSecret, notificationapi.NotificationType(openapi.TestMessage)); err != nil {
		if errors.Is(err, notification.ErrNoNotificationChannels) {
			return nil, err
		}
		return nil, fmt.Errorf("cannot send test notification: %w", err)
	}

	return openapi.SendTestNotification204Response{}, nil
}
