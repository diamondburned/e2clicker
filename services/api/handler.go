package api

import (
	"context"
	"log/slog"

	"libdb.so/e2clicker/services/api/openapi"
	"libdb.so/e2clicker/services/user"
)

type oapiHandler struct {
	users  *user.UserService `do:""`
	logger *slog.Logger      `do:""`
}

// (POST /login)
func (h *oapiHandler) Login(ctx context.Context, request openapi.LoginRequestObject) (openapi.LoginResponseObject, error) {
	uID, err := h.users.ValidateUserEmailPassword(ctx, request.Body.Email, request.Body.Password)
	if err != nil {
		return convertError[openapi.LogindefaultJSONResponse](ctx, err), nil
	}

	token, err := h.users.RegisterSession(ctx, uID, optstr(request.Params.UserAgent))
	if err != nil {
		return convertError[openapi.LogindefaultJSONResponse](ctx, err), nil
	}

	return openapi.Login200JSONResponse{
		UserID: uID,
		Token:  token,
	}, nil
}

// (POST /register)
func (h *oapiHandler) Register(ctx context.Context, request openapi.RegisterRequestObject) (openapi.RegisterResponseObject, error) {
	u, err := h.users.CreateUser(ctx, request.Body.Email, request.Body.Password, request.Body.Name)
	if err != nil {
		return convertError[openapi.RegisterdefaultJSONResponse](ctx, err), nil
	}

	token, err := h.users.RegisterSession(ctx, u.ID, optstr(request.Params.UserAgent))
	if err != nil {
		return convertError[openapi.RegisterdefaultJSONResponse](ctx, err), nil
	}

	return openapi.Register200JSONResponse{
		User:  convertUser(u),
		Token: token,
	}, nil
}

// (GET /user/{userID})
func (h *oapiHandler) User(ctx context.Context, request openapi.UserRequestObject) (openapi.UserResponseObject, error) {
	u, err := h.users.User(ctx, request.UserIDParam)
	if err != nil {
		return convertError[openapi.UserdefaultJSONResponse](ctx, err), nil
	}

	return openapi.User200JSONResponse(convertUser(u)), nil
}

// (GET /user/{userID}/avatar)
func (h *oapiHandler) UserAvatar(ctx context.Context, request openapi.UserAvatarRequestObject) (openapi.UserAvatarResponseObject, error) {
	u, err := h.users.UserAvatar(ctx, request.UserIDParam)
	if err != nil {
		return convertError[openapi.UserAvatardefaultJSONResponse](ctx, err), nil
	}

	return openapi.UserAvatar200ImageResponse{
		Body:          u.Reader(),
		ContentType:   u.ContentType,
		ContentLength: u.ContentLength,
	}, nil
}
