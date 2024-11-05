package api

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
	"libdb.so/e2clicker/internal/asset"
	"libdb.so/e2clicker/services/api/openapi"
	"libdb.so/e2clicker/services/dosage"
	"libdb.so/e2clicker/services/user"
)

// OpenAPIHandler is the handler for the OpenAPI service.
// It implements the OpenAPI service interface.
type OpenAPIHandler struct {
	logger *slog.Logger
	users  *user.UserService
	dosage dosage.DosageStorage
}

// OpenAPIHandlerServices is the set of service dependencies required by the
// OpenAPIHandler.
type OpenAPIHandlerServices struct {
	fx.In

	*user.UserService
	dosage.DosageStorage
}

// NewOpenAPIHandler creates a new OpenAPIHandler.
func NewOpenAPIHandler(deps OpenAPIHandlerServices, logger *slog.Logger) *OpenAPIHandler {
	return &OpenAPIHandler{
		logger: logger,
		users:  deps.UserService,
		dosage: deps.DosageStorage,
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

	return openapi.CurrentUser200JSONResponse(convertUser(u)), nil
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

// Get the current user's secret
// (GET /me/secret)
func (h *OpenAPIHandler) CurrentUserSecret(ctx context.Context, request openapi.CurrentUserSecretRequestObject) (openapi.CurrentUserSecretResponseObject, error) {
	session := sessionFromCtx(ctx)
	return openapi.CurrentUserSecret200JSONResponse{
		Secret: session.UserSecret,
	}, nil
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

	err := h.users.DeleteSession(ctx, session.UserSecret, request.Body.ID)
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

func (h *OpenAPIHandler) DoseHistory(ctx context.Context, request openapi.DoseHistoryRequestObject) (openapi.DoseHistoryResponseObject, error) {
	session := sessionFromCtx(ctx)

	s, err := h.dosage.DoseHistory(ctx, session.UserSecret, request.Body.Start, request.Body.End)
	if err != nil {
		return nil, err
	}

	return openapi.DoseHistory200JSONResponse(
		convertList(s.Entries, convertDosageObservation),
	), nil
}

func (h *OpenAPIHandler) RecordDose(ctx context.Context, request openapi.RecordDoseRequestObject) (openapi.RecordDoseResponseObject, error) {
	session := sessionFromCtx(ctx)

	o, err := h.dosage.RecordDose(ctx, session.UserSecret, request.Body.TakenAt)
	if err != nil {
		return nil, err
	}

	return openapi.RecordDose200JSONResponse(convertDosageObservation(o)), nil
}

func (h *OpenAPIHandler) EditDose(ctx context.Context, request openapi.EditDoseRequestObject) (openapi.EditDoseResponseObject, error) {
	session := sessionFromCtx(ctx)

	o := dosage.Observation{
		DoseID:         request.Body.ID,
		DeliveryMethod: request.Body.DeliveryMethod,
		Dose:           request.Body.Dose,
		TakenAt:        request.Body.TakenAt,
		TakenOffAt:     request.Body.TakenOffAt,
	}

	if err := h.dosage.EditDose(ctx, session.UserSecret, o); err != nil {
		return nil, err
	}

	return openapi.EditDose204Response{}, nil
}

func (h *OpenAPIHandler) ForgetDoses(ctx context.Context, request openapi.ForgetDosesRequestObject) (openapi.ForgetDosesResponseObject, error) {
	session := sessionFromCtx(ctx)

	if err := h.dosage.ForgetDoses(ctx, session.UserSecret, request.Body.DoseIDs); err != nil {
		return nil, err
	}

	return openapi.ForgetDoses204Response{}, nil
}

func convertDosageObservation(o dosage.Observation) openapi.DosageObservation {
	return openapi.DosageObservation{
		ID:             o.DoseID,
		DeliveryMethod: o.DeliveryMethod,
		Dose:           o.Dose,
		TakenAt:        o.TakenAt,
		TakenOffAt:     o.TakenOffAt,
	}
}

// Get the user's dosage schedule
// (GET /dosage/schedule)
func (h *OpenAPIHandler) DosageSchedule(ctx context.Context, request openapi.DosageScheduleRequestObject) (openapi.DosageScheduleResponseObject, error) {
	session := sessionFromCtx(ctx)

	s, err := h.dosage.DosageSchedule(ctx, session.UserSecret)
	if err != nil {
		return nil, err
	}

	return openapi.DosageSchedule200JSONResponse(openapi.DosageSchedule{
		DeliveryMethod: s.DeliveryMethod,
		Dose:           s.Dose,
		Interval:       float64(s.Interval),
		Concurrence:    s.Concurrence,
	}), nil
}

// Set the user's dosage schedule
// (PUT /dosage/schedule)
func (h *OpenAPIHandler) SetDosageSchedule(ctx context.Context, request openapi.SetDosageScheduleRequestObject) (openapi.SetDosageScheduleResponseObject, error) {
	session := sessionFromCtx(ctx)

	if err := h.dosage.SetDosageSchedule(ctx, dosage.Schedule{
		UserSecret:     session.UserSecret,
		DeliveryMethod: request.Body.DeliveryMethod,
		Dose:           request.Body.Dose,
		Interval:       dosage.Days(request.Body.Interval),
		Concurrence:    request.Body.Concurrence,
	}); err != nil {
		return nil, err
	}

	return openapi.SetDosageSchedule204Response{}, nil
}
