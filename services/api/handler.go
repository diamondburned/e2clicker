package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"e2clicker.app/internal/publicerrors"
	"e2clicker.app/services/api/openapi"
	"e2clicker.app/services/dosage"
	"e2clicker.app/services/notification"
	"e2clicker.app/services/user"
	"go.uber.org/fx"

	notificationapi "e2clicker.app/services/notification/openapi"
)

// openAPIHandler is the handler for the OpenAPI service.
// It implements the OpenAPI service interface.
type openAPIHandler struct {
	logger      *slog.Logger
	users       *user.UserService
	notifs      *notification.UserNotificationService
	notif       *notification.NotificationService
	dosage      dosage.DosageStorage
	doseHistory dosage.DoseHistoryStorage
}

// OpenAPIHandlerServices is the set of service dependencies required by the
// OpenAPIHandler.
type OpenAPIHandlerServices struct {
	fx.In

	Users             *user.UserService
	UserNotifications *notification.UserNotificationService
	Notification      *notification.NotificationService
	Dosage            dosage.DosageStorage
	DoseHistory       dosage.DoseHistoryStorage
}

// newOpenAPIHandler creates a new OpenAPIHandler.
func newOpenAPIHandler(deps OpenAPIHandlerServices, logger *slog.Logger) *openAPIHandler {
	return &openAPIHandler{
		logger:      logger,
		users:       deps.Users,
		notifs:      deps.UserNotifications,
		notif:       deps.Notification,
		dosage:      deps.Dosage,
		doseHistory: deps.DoseHistory,
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

	s := dosage.Dosage{
		UserSecret:     session.UserSecret,
		DeliveryMethod: request.Body.DeliveryMethod,
		Dose:           request.Body.Dose,
		Interval:       dosage.Days(request.Body.Interval),
		Concurrence:    request.Body.Concurrence,
	}

	if !slices.ContainsFunc(methods, func(m dosage.DeliveryMethod) bool {
		return m.ID == s.DeliveryMethod
	}) {
		return nil, publicerrors.Errorf("invalid delivery method %q", s.DeliveryMethod)
	}

	if err := h.dosage.SetDosage(ctx, s); err != nil {
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
	now := time.Now()

	d, err := h.dosage.Dosage(ctx, session.UserSecret)
	if err != nil {
		return nil, publicerrors.New("no dosage set")
	}

	dose := dosage.Dose{
		DeliveryMethod: d.DeliveryMethod,
		Dose:           d.Dose,
		TakenAt:        now,
	}

	if err := h.doseHistory.RecordDose(ctx, session.UserSecret, dose); err != nil {
		return nil, err
	}

	return openapi.RecordDose200JSONResponse(openapi.Dose(dose.ToOpenAPI())), nil
}

func (h *openAPIHandler) EditDose(ctx context.Context, request openapi.EditDoseRequestObject) (openapi.EditDoseResponseObject, error) {
	session := sessionFromCtx(ctx)

	o := dosage.Dose{
		DeliveryMethod: request.Body.DeliveryMethod,
		Dose:           request.Body.Dose,
		TakenAt:        request.Body.TakenAt,
		TakenOffAt:     request.Body.TakenOffAt,
	}

	if err := h.doseHistory.EditDose(ctx, session.UserSecret, request.DoseTime, o); err != nil {
		return nil, err
	}

	return openapi.EditDose204Response{}, nil
}

func (h *openAPIHandler) ForgetDose(ctx context.Context, request openapi.ForgetDoseRequestObject) (openapi.ForgetDoseResponseObject, error) {
	session := sessionFromCtx(ctx)
	if err := h.doseHistory.ForgetDoses(ctx, session.UserSecret, []time.Time{request.DoseTime}); err != nil {
		return nil, err
	}
	return openapi.ForgetDose204Response{}, nil
}

func (h *openAPIHandler) ForgetDoses(ctx context.Context, request openapi.ForgetDosesRequestObject) (openapi.ForgetDosesResponseObject, error) {
	session := sessionFromCtx(ctx)
	if err := h.doseHistory.ForgetDoses(ctx, session.UserSecret, request.Params.DoseTimes); err != nil {
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
		r.Dosage = &openapi.Dosage{
			DeliveryMethod: dosage.DeliveryMethod,
			Dose:           dosage.Dose,
			Interval:       float64(dosage.Interval),
			Concurrence:    dosage.Concurrence,
		}
	}

	if request.Params.Start != nil && request.Params.End != nil {
		const oneYear = 365 * 24 * time.Hour
		if request.Params.End.Sub(*request.Params.Start) > oneYear {
			return nil, publicerrors.New("" +
				"requested history range is too large, must be 1 year or less " +
				"(consider exporting instead)")
		}

		os := make([]openapi.Dose, 0, 32)
		r.History = &os

		for dose, err := range h.doseHistory.DoseHistory(
			ctx, session.UserSecret,
			*request.Params.Start,
			*request.Params.End) {

			if err != nil {
				return nil, fmt.Errorf("cannot get dosage history: %w", err)
			}

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
	i, err := h.notifs.WebPushInfo(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.WebPushInfo200JSONResponse(openapi.PushInfo(i)), nil
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

	supports := h.notif.Supports()

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
		s := make([]openapi.PushSubscription, len(p.NotificationConfigs.WebPush))
		for i, sub := range p.NotificationConfigs.WebPush {
			s[i] = openapi.PushSubscription{
				DeviceID:       sub.DeviceID,
				ExpirationTime: sub.ExpirationTime,
			}
			s[i].Keys.P256Dh = sub.Keys.P256Dh
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

	if request.Body.NotificationConfigs.WebPush != nil {
		newPreferences.NotificationConfigs.WebPush = make([]notificationapi.PushSubscription, len(*request.Body.NotificationConfigs.WebPush))
		for i, v := range *request.Body.NotificationConfigs.WebPush {
			newPreferences.NotificationConfigs.WebPush[i] = notificationapi.PushSubscription(v)
		}
	}

	if err := h.notifs.SetUserPreferences(ctx, session.UserSecret, newPreferences); err != nil {
		return nil, err
	}

	return openapi.UserUpdateNotificationPreferences204Response{}, nil
}

// Send a test notification
// (POST /notifications/test)
func (h *openAPIHandler) SendTestNotification(ctx context.Context, request openapi.SendTestNotificationRequestObject) (openapi.SendTestNotificationResponseObject, error) {
	session := sessionFromCtx(ctx)

	if err := h.notifs.NotifyUser(ctx, session.UserSecret, notificationapi.NotificationType(openapi.TestMessage)); err != nil {
		return nil, fmt.Errorf("cannot send test notification: %w", err)
	}

	return openapi.SendTestNotification204Response{}, nil
}
