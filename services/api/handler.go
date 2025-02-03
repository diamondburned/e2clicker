package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"go.uber.org/fx"
	"e2clicker.app/internal/publicerrors"
	"e2clicker.app/services/api/openapi"
	"e2clicker.app/services/dosage"
	"e2clicker.app/services/notification"
	"e2clicker.app/services/user"

	notificationapi "e2clicker.app/services/notification/openapi"
)

// openAPIHandler is the handler for the OpenAPI service.
// It implements the OpenAPI service interface.
type openAPIHandler struct {
	logger      *slog.Logger
	users       *user.UserService
	notifs      *notification.UserNotificationService
	dosage      dosage.DosageStorage
	doseHistory dosage.DoseHistoryStorage
}

// OpenAPIHandlerServices is the set of service dependencies required by the
// OpenAPIHandler.
type OpenAPIHandlerServices struct {
	fx.In

	Users             *user.UserService
	UserNotifications *notification.UserNotificationService
	Dosage            dosage.DosageStorage
	DoseHistory       dosage.DoseHistoryStorage
}

// newOpenAPIHandler creates a new OpenAPIHandler.
func newOpenAPIHandler(deps OpenAPIHandlerServices, logger *slog.Logger) *openAPIHandler {
	return &openAPIHandler{
		logger:      logger,
		users:       deps.Users,
		notifs:      deps.UserNotifications,
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
		Name:      u.Name,
		Locale:    u.Locale,
		Secret:    u.Secret,
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

func (h *openAPIHandler) UserNotificationMethods(ctx context.Context, request openapi.UserNotificationMethodsRequestObject) (openapi.UserNotificationMethodsResponseObject, error) {
	session := sessionFromCtx(ctx)

	p, err := h.notifs.UserPreferences(ctx, session.UserSecret)
	if err != nil {
		return nil, err
	}

	var ret openapi.ReturnedNotificationMethods

	if len(p.NotificationConfigs.WebPush) > 0 {
		s := make([]openapi.ReturnedPushSubscription, len(p.NotificationConfigs.WebPush))
		for i, sub := range p.NotificationConfigs.WebPush {
			s[i] = openapi.ReturnedPushSubscription{
				DeviceID:       sub.DeviceID,
				ExpirationTime: sub.ExpirationTime,
			}
			s[i].Keys.P256Dh = sub.Keys.P256Dh
		}
		ret.WebPush = &s
	}

	return openapi.UserNotificationMethods200JSONResponse(ret), nil
}

func (h *openAPIHandler) UserPushSubscription(ctx context.Context, request openapi.UserPushSubscriptionRequestObject) (openapi.UserPushSubscriptionResponseObject, error) {
	session := sessionFromCtx(ctx)

	p, err := h.notifs.UserPreferences(ctx, session.UserSecret)
	if err != nil {
		return nil, err
	}

	ix := slices.IndexFunc(p.NotificationConfigs.WebPush,
		func(c notificationapi.PushSubscription) bool { return c.DeviceID == request.DeviceID },
	)
	if ix == -1 {
		return openapi.UserPushSubscription404JSONResponse(openapi.Error{
			Message: "subscription not found",
			Details: anyPtr(map[string]string{
				"deviceID": string(request.DeviceID),
			}),
		}), nil
	}

	return openapi.UserPushSubscription200JSONResponse(openapi.PushSubscription(
		p.NotificationConfigs.WebPush[ix],
	)), nil
}

func (h *openAPIHandler) UserSubscribePush(ctx context.Context, request openapi.UserSubscribePushRequestObject) (openapi.UserSubscribePushResponseObject, error) {
	session := sessionFromCtx(ctx)
	if err := h.notifs.SubscribeWebPush(ctx, session.UserSecret, notificationapi.PushSubscription(*request.Body)); err != nil {
		return nil, err
	}
	return openapi.UserSubscribePush204Response{}, nil
}

func (h *openAPIHandler) UserUnsubscribePush(ctx context.Context, request openapi.UserUnsubscribePushRequestObject) (openapi.UserUnsubscribePushResponseObject, error) {
	session := sessionFromCtx(ctx)
	if err := h.notifs.UnsubscribeWebPush(ctx, session.UserSecret, request.DeviceID); err != nil {
		return nil, err
	}
	return openapi.UserUnsubscribePush204Response{}, nil
}
