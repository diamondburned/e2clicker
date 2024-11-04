package api

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
	"libdb.so/e2clicker/internal/asset"
	"libdb.so/e2clicker/services/api/openapi"
	"libdb.so/e2clicker/services/user"
)

// OpenAPIHandler is the handler for the OpenAPI service.
// It implements the OpenAPI service interface.
type OpenAPIHandler struct {
	logger *slog.Logger
	users  *user.UserService
}

// OpenAPIHandlerServices is the set of service dependencies required by the
// OpenAPIHandler.
type OpenAPIHandlerServices struct {
	fx.In

	*user.UserService
}

// NewOpenAPIHandler creates a new OpenAPIHandler.
func NewOpenAPIHandler(deps OpenAPIHandlerServices, logger *slog.Logger) *OpenAPIHandler {
	return &OpenAPIHandler{
		logger: logger,
		users:  deps.UserService,
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
