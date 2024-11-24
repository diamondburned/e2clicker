package api

import (
	"context"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"slices"
	"time"

	"github.com/timewasted/go-accept-headers"
	"go.uber.org/fx"
	"libdb.so/e2clicker/internal/asset"
	"libdb.so/e2clicker/internal/publicerrors"
	"libdb.so/e2clicker/services/api/openapi"
	"libdb.so/e2clicker/services/dosage"
	"libdb.so/e2clicker/services/notification"
	"libdb.so/e2clicker/services/user"

	notificationapi "libdb.so/e2clicker/services/notification/openapi"
)

// OpenAPIHandler is the handler for the OpenAPI service.
// It implements the OpenAPI service interface.
type OpenAPIHandler struct {
	logger       *slog.Logger
	users        *user.UserService
	notifs       *notification.UserNotificationService
	dosage       dosage.DosageStorage
	doseHistory  dosage.DoseHistoryStorage
	doseExporter *dosage.ExporterService
}

// OpenAPIHandlerServices is the set of service dependencies required by the
// OpenAPIHandler.
type OpenAPIHandlerServices struct {
	fx.In

	Users             *user.UserService
	UserNotifications *notification.UserNotificationService
	Dosage            dosage.DosageStorage
	DoseHistory       dosage.DoseHistoryStorage
	DoseExporter      *dosage.ExporterService
}

// NewOpenAPIHandler creates a new OpenAPIHandler.
func NewOpenAPIHandler(deps OpenAPIHandlerServices, logger *slog.Logger) *OpenAPIHandler {
	return &OpenAPIHandler{
		logger:       logger,
		users:        deps.Users,
		notifs:       deps.UserNotifications,
		dosage:       deps.Dosage,
		doseHistory:  deps.DoseHistory,
		doseExporter: deps.DoseExporter,
	}
}

func (h *OpenAPIHandler) asStrictHandler() openapi.StrictServerInterface { return h }

// Register a new account
// (POST /register)
func (h *OpenAPIHandler) Register(ctx context.Context, request openapi.RegisterRequestObject) (openapi.RegisterResponseObject, error) {
	u, err := h.users.CreateUser(ctx, request.Body.Name)
	if err != nil {
		return nil, err
	}

	return openapi.Register200JSONResponse{
		Name:      u.Name,
		Locale:    u.Locale,
		HasAvatar: u.HasAvatar,
		Secret:    u.Secret,
	}, nil
}

// Authenticate a user
// (POST /auth)
func (h *OpenAPIHandler) Auth(ctx context.Context, request openapi.AuthRequestObject) (openapi.AuthResponseObject, error) {
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
func (h *OpenAPIHandler) CurrentUser(ctx context.Context, request openapi.CurrentUserRequestObject) (openapi.CurrentUserResponseObject, error) {
	session := sessionFromCtx(ctx)

	u, err := h.users.User(ctx, session.UserSecret)
	if err != nil {
		return nil, err
	}

	return openapi.CurrentUser200JSONResponse{
		Name:      u.Name,
		Locale:    u.Locale,
		HasAvatar: u.HasAvatar,
		Secret:    session.UserSecret,
	}, nil
}

// Get the current user's avatar
// (GET /me/avatar)
func (h *OpenAPIHandler) CurrentUserAvatar(ctx context.Context, request openapi.CurrentUserAvatarRequestObject) (openapi.CurrentUserAvatarResponseObject, error) {
	session := sessionFromCtx(ctx)

	a, err := h.users.UserAvatar(ctx, session.UserSecret)
	if err != nil {
		return nil, err
	}

	return openapi.CurrentUserAvatar200ImageResponse{
		Body:          a.Reader(),
		ContentType:   a.ContentType,
		ContentLength: a.ContentLength,
	}, nil
}

// Set the current user's avatar
// (POST /me/avatar)
func (h *OpenAPIHandler) SetCurrentUserAvatar(ctx context.Context, request openapi.SetCurrentUserAvatarRequestObject) (openapi.SetCurrentUserAvatarResponseObject, error) {
	session := sessionFromCtx(ctx)

	err := h.users.SetUserAvatar(ctx, session.UserSecret, asset.NewAssetReader(
		request.Body,
		request.ContentType,
		-1,
	))
	if err != nil {
		return nil, err
	}

	return openapi.SetCurrentUserAvatar204Response{}, nil
}

// List the current user's sessions
// (GET /me/sessions)
func (h *OpenAPIHandler) CurrentUserSessions(ctx context.Context, request openapi.CurrentUserSessionsRequestObject) (openapi.CurrentUserSessionsResponseObject, error) {
	session := sessionFromCtx(ctx)

	s, err := h.users.ListSessions(ctx, session.UserSecret)
	if err != nil {
		return nil, err
	}

	return openapi.CurrentUserSessions200JSONResponse(convertList(s, convertSession)), nil
}

// Delete one of the current user's sessions
// (DELETE /me/sessions)
func (h *OpenAPIHandler) DeleteUserSession(ctx context.Context, request openapi.DeleteUserSessionRequestObject) (openapi.DeleteUserSessionResponseObject, error) {
	session := sessionFromCtx(ctx)

	err := h.users.DeleteSession(ctx, session.UserSecret, request.Params.ID)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteUserSession204Response{}, nil
}

// List all available delivery methods
// (GET /delivery-methods)
func (h *OpenAPIHandler) DeliveryMethods(ctx context.Context, request openapi.DeliveryMethodsRequestObject) (openapi.DeliveryMethodsResponseObject, error) {
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

func (h *OpenAPIHandler) SetDosage(ctx context.Context, request openapi.SetDosageRequestObject) (openapi.SetDosageResponseObject, error) {
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

func (h *OpenAPIHandler) ClearDosage(ctx context.Context, request openapi.ClearDosageRequestObject) (openapi.ClearDosageResponseObject, error) {
	session := sessionFromCtx(ctx)
	if err := h.dosage.ClearDosage(ctx, session.UserSecret); err != nil {
		return nil, err
	}
	return openapi.ClearDosage204Response{}, nil
}

func (h *OpenAPIHandler) RecordDose(ctx context.Context, request openapi.RecordDoseRequestObject) (openapi.RecordDoseResponseObject, error) {
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

	id, err := h.doseHistory.RecordDose(ctx, session.UserSecret, dose)
	if err != nil {
		return nil, err
	}

	return openapi.RecordDose200JSONResponse(
		convertDosageObservation(dosage.Observation{
			ID:   id,
			Dose: dose,
		}),
	), nil
}

func (h *OpenAPIHandler) EditDose(ctx context.Context, request openapi.EditDoseRequestObject) (openapi.EditDoseResponseObject, error) {
	session := sessionFromCtx(ctx)

	o := dosage.Dose{
		DeliveryMethod: request.Body.DeliveryMethod,
		Dose:           request.Body.Dose,
		TakenAt:        request.Body.TakenAt,
		TakenOffAt:     request.Body.TakenOffAt,
	}

	if err := h.doseHistory.EditDose(ctx, session.UserSecret, request.Body.ID, o); err != nil {
		return nil, err
	}

	return openapi.EditDose204Response{}, nil
}

func (h *OpenAPIHandler) ForgetDoses(ctx context.Context, request openapi.ForgetDosesRequestObject) (openapi.ForgetDosesResponseObject, error) {
	session := sessionFromCtx(ctx)

	if err := h.doseHistory.ForgetDoses(ctx, session.UserSecret, request.Params.DoseIds); err != nil {
		return nil, err
	}

	return openapi.ForgetDoses204Response{}, nil
}

func (h *OpenAPIHandler) Dosage(ctx context.Context, request openapi.DosageRequestObject) (openapi.DosageResponseObject, error) {
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

		os := make([]openapi.DosageObservation, 0, 32)
		r.History = &os

		for dose, err := range h.doseHistory.DoseHistory(
			ctx, session.UserSecret,
			*request.Params.Start,
			*request.Params.End) {

			if err != nil {
				return nil, fmt.Errorf("cannot get dosage history: %w", err)
			}

			os = append(os, convertDosageObservation(dose))
		}
	}

	return r, nil
}

func (h *OpenAPIHandler) ExportDosageHistory(ctx context.Context, request openapi.ExportDosageHistoryRequestObject) (openapi.ExportDosageHistoryResponseObject, error) {
	session := sessionFromCtx(ctx)

	var format dosage.ExportFormat
	for _, t := range accept.Parse(string(request.Params.Accept)) {
		switch t.Type + "/" + t.Subtype {
		case "text/csv":
			format = dosage.ExportCSV
		case "application/json":
			format = dosage.ExportJSON
		}
	}

	if format == "" {
		return nil, ErrNoAcceptableContentType
	}

	exportExtensions, err := mime.ExtensionsByType(string(format))
	if err != nil {
		return nil, fmt.Errorf("format %q missing file extension: %w", format, err)
	}

	exportTime := time.Now().Format(time.RFC3339)
	exportName := fmt.Sprintf("attachment; filename=dose-history-%s.%s", exportTime, exportExtensions[0])

	return exportDosageHistoryResponse(func(w http.ResponseWriter) error {
		w.Header().Set("Content-Type", string(request.Params.Accept))
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", exportName))

		_, err := h.doseExporter.ExportDoseHistory(ctx, w, session.UserSecret, dosage.ExportDoseHistoryOptions{
			Begin:  optPtr(request.Params.Start),
			End:    optPtr(request.Params.End),
			Format: dosage.ExportFormat(request.Params.Accept),
		})
		return err
	}), nil
}

type exportDosageHistoryResponse func(http.ResponseWriter) error

func (f exportDosageHistoryResponse) VisitExportDosageHistoryResponse(w http.ResponseWriter) error {
	return f(w)
}

func (h *OpenAPIHandler) ImportDosageHistory(ctx context.Context, request openapi.ImportDosageHistoryRequestObject) (openapi.ImportDosageHistoryResponseObject, error) {
	session := sessionFromCtx(ctx)

	contentType, params, err := mime.ParseMediaType(string(request.Params.ContentType))
	if err != nil {
		return nil, publicerrors.Errorf("invalid content type: %w", err)
	}

	if charset, ok := params["charset"]; ok && charset != "utf-8" {
		return nil, publicerrors.Errorf("unsupported charset %q, UTF-8 only please", charset)
	}

	var format dosage.ExportFormat
	switch contentType {
	case "text/csv":
		format = dosage.ExportCSV
	case "application/json":
		format = dosage.ExportJSON
	default:
		return nil, ErrNoAcceptableContentType
	}

	result, err := h.doseExporter.ImportDoseHistory(ctx, request.Body, session.UserSecret, dosage.ImportDoseHistoryOptions{
		Format: format,
	})
	if result.Records == 0 && err != nil {
		return nil, err
	}

	var oapiError *openapi.Error
	if err != nil {
		converted := convertError[errorResponse](ctx, err)
		oapiError = &converted.Body
	}

	return openapi.ImportDosageHistory200JSONResponse{
		Records:   int(result.Records),
		Succeeded: int(result.Succeeded),
		Error:     oapiError,
	}, nil
}

func convertDosageObservation(o dosage.Observation) openapi.DosageObservation {
	return openapi.DosageObservation{
		ID:             o.ID,
		DeliveryMethod: o.Dose.DeliveryMethod,
		Dose:           o.Dose.Dose,
		TakenAt:        o.Dose.TakenAt,
		TakenOffAt:     o.Dose.TakenOffAt,
	}
}

func (h *OpenAPIHandler) WebPushInfo(ctx context.Context, request openapi.WebPushInfoRequestObject) (openapi.WebPushInfoResponseObject, error) {
	i, err := h.notifs.WebPushInfo(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.WebPushInfo200JSONResponse(openapi.PushInfo(i)), nil
}

func (h *OpenAPIHandler) UserNotificationMethods(ctx context.Context, request openapi.UserNotificationMethodsRequestObject) (openapi.UserNotificationMethodsResponseObject, error) {
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

func (h *OpenAPIHandler) UserPushSubscription(ctx context.Context, request openapi.UserPushSubscriptionRequestObject) (openapi.UserPushSubscriptionResponseObject, error) {
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

func (h *OpenAPIHandler) UserSubscribePush(ctx context.Context, request openapi.UserSubscribePushRequestObject) (openapi.UserSubscribePushResponseObject, error) {
	session := sessionFromCtx(ctx)
	if err := h.notifs.SubscribeWebPush(ctx, session.UserSecret, notificationapi.PushSubscription(*request.Body)); err != nil {
		return nil, err
	}
	return openapi.UserSubscribePush204Response{}, nil
}

func (h *OpenAPIHandler) UserUnsubscribePush(ctx context.Context, request openapi.UserUnsubscribePushRequestObject) (openapi.UserUnsubscribePushResponseObject, error) {
	session := sessionFromCtx(ctx)
	if err := h.notifs.UnsubscribeWebPush(ctx, session.UserSecret, request.DeviceID); err != nil {
		return nil, err
	}
	return openapi.UserUnsubscribePush204Response{}, nil
}
