// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version (devel) DO NOT EDIT.
package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"libdb.so/e2clicker/services/user"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Error defines model for Error.
type Error struct {
	// Message A message describing the error
	Message string `json:"message"`

	// Details Additional details about the error
	Details *interface{} `json:"details,omitempty"`

	// Internal Whether the error is internal
	Internal *bool `json:"internal,omitempty"`

	// InternalCode An internal code for the error (useless for clients)
	InternalCode *string `json:"internalCode,omitempty"`
}

// Locale A locale identifier.
type Locale = user.Locale

// SessionToken A session token string. This is used in the Authorization header to authenticate requests.
type SessionToken = user.SessionToken

// User A user of the system.
type User struct {
	// ID A unique user identifier.
	ID UserID `json:"id"`

	// Email The user's email address
	Email string `json:"email"`

	// Name The user's name
	Name string `json:"name"`

	// Locale A locale identifier.
	Locale Locale `json:"locale"`
}

// UserID A unique user identifier.
type UserID = user.UserID

// UserIDParam defines model for userID.
type UserIDParam = user.UserID

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse = Error

// LoginJSONBody defines parameters for Login.
type LoginJSONBody struct {
	// Email The username to log in with
	Email string `json:"email"`

	// Password The password to log in with
	Password string `json:"password"`
}

// LoginParams defines parameters for Login.
type LoginParams struct {
	// UserAgent The user agent string of the client making the request.
	UserAgent *string `json:"User-Agent,omitempty"`
}

// RegisterJSONBody defines parameters for Register.
type RegisterJSONBody struct {
	// Name The name to register with
	Name string `json:"name"`

	// Email The username to register with
	Email string `json:"email"`

	// Password The password to register with
	Password string `json:"password"`
}

// RegisterParams defines parameters for Register.
type RegisterParams struct {
	// UserAgent The user agent string of the client making the request.
	UserAgent *string `json:"User-Agent,omitempty"`
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody LoginJSONBody

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody RegisterJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Log into an existing account
	// (POST /login)
	Login(w http.ResponseWriter, r *http.Request, params LoginParams)
	// Register a new account
	// (POST /register)
	Register(w http.ResponseWriter, r *http.Request, params RegisterParams)
	// Get a user by ID
	// (GET /user/{userID})
	User(w http.ResponseWriter, r *http.Request, userIDParam UserIDParam)
	// Get a user's avatar by ID
	// (GET /user/{userID}/avatar)
	UserAvatar(w http.ResponseWriter, r *http.Request, userIDParam UserIDParam)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Log into an existing account
// (POST /login)
func (_ Unimplemented) Login(w http.ResponseWriter, r *http.Request, params LoginParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Register a new account
// (POST /register)
func (_ Unimplemented) Register(w http.ResponseWriter, r *http.Request, params RegisterParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get a user by ID
// (GET /user/{userID})
func (_ Unimplemented) User(w http.ResponseWriter, r *http.Request, userIDParam UserIDParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get a user's avatar by ID
// (GET /user/{userID}/avatar)
func (_ Unimplemented) UserAvatar(w http.ResponseWriter, r *http.Request, userIDParam UserIDParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// Login operation middleware
func (siw *ServerInterfaceWrapper) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params LoginParams

	headers := r.Header

	// ------------- Optional header parameter "User-Agent" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("User-Agent")]; found {
		var UserAgent string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "User-Agent", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "User-Agent", valueList[0], &UserAgent, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "User-Agent", Err: err})
			return
		}

		params.UserAgent = &UserAgent

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Login(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Register operation middleware
func (siw *ServerInterfaceWrapper) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params RegisterParams

	headers := r.Header

	// ------------- Optional header parameter "User-Agent" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("User-Agent")]; found {
		var UserAgent string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "User-Agent", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "User-Agent", valueList[0], &UserAgent, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "User-Agent", Err: err})
			return
		}

		params.UserAgent = &UserAgent

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Register(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// User operation middleware
func (siw *ServerInterfaceWrapper) User(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "userID" -------------
	var userIDParam UserIDParam

	err = runtime.BindStyledParameterWithOptions("simple", "userID", chi.URLParam(r, "userID"), &userIDParam, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "userID", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.User(w, r, userIDParam)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// UserAvatar operation middleware
func (siw *ServerInterfaceWrapper) UserAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "userID" -------------
	var userIDParam UserIDParam

	err = runtime.BindStyledParameterWithOptions("simple", "userID", chi.URLParam(r, "userID"), &userIDParam, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "userID", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UserAvatar(w, r, userIDParam)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/login", wrapper.Login)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/register", wrapper.Register)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/user/{userID}", wrapper.User)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/user/{userID}/avatar", wrapper.UserAvatar)
	})

	return r
}

type ErrorResponseJSONResponse Error

type LoginRequestObject struct {
	Params LoginParams
	Body   *LoginJSONRequestBody
}

type LoginResponseObject interface {
	VisitLoginResponse(w http.ResponseWriter) error
}

type Login200JSONResponse struct {
	// UserID A unique user identifier.
	UserID UserID `json:"userID"`

	// Token A session token string. This is used in the Authorization header to authenticate requests.
	Token SessionToken `json:"token"`
}

func (response Login200JSONResponse) VisitLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type LogindefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response LogindefaultJSONResponse) VisitLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type RegisterRequestObject struct {
	Params RegisterParams
	Body   *RegisterJSONRequestBody
}

type RegisterResponseObject interface {
	VisitRegisterResponse(w http.ResponseWriter) error
}

type Register200JSONResponse struct {
	// User A user of the system.
	User User `json:"user"`

	// Token A session token string. This is used in the Authorization header to authenticate requests.
	Token SessionToken `json:"token"`
}

func (response Register200JSONResponse) VisitRegisterResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type RegisterdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response RegisterdefaultJSONResponse) VisitRegisterResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type UserRequestObject struct {
	UserIDParam UserIDParam `json:"userID"`
}

type UserResponseObject interface {
	VisitUserResponse(w http.ResponseWriter) error
}

type User200JSONResponse User

func (response User200JSONResponse) VisitUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type UserdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response UserdefaultJSONResponse) VisitUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type UserAvatarRequestObject struct {
	UserIDParam UserIDParam `json:"userID"`
}

type UserAvatarResponseObject interface {
	VisitUserAvatarResponse(w http.ResponseWriter) error
}

type UserAvatar200ImageResponse struct {
	Body          io.Reader
	ContentType   string
	ContentLength int64
}

func (response UserAvatar200ImageResponse) VisitUserAvatarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", response.ContentType)
	if response.ContentLength != 0 {
		w.Header().Set("Content-Length", fmt.Sprint(response.ContentLength))
	}
	w.WriteHeader(200)

	if closer, ok := response.Body.(io.ReadCloser); ok {
		defer closer.Close()
	}
	_, err := io.Copy(w, response.Body)
	return err
}

type UserAvatardefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response UserAvatardefaultJSONResponse) VisitUserAvatarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Log into an existing account
	// (POST /login)
	Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error)
	// Register a new account
	// (POST /register)
	Register(ctx context.Context, request RegisterRequestObject) (RegisterResponseObject, error)
	// Get a user by ID
	// (GET /user/{userID})
	User(ctx context.Context, request UserRequestObject) (UserResponseObject, error)
	// Get a user's avatar by ID
	// (GET /user/{userID}/avatar)
	UserAvatar(ctx context.Context, request UserAvatarRequestObject) (UserAvatarResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// Login operation middleware
func (sh *strictHandler) Login(w http.ResponseWriter, r *http.Request, params LoginParams) {
	var request LoginRequestObject

	request.Params = params

	var body LoginJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Login(ctx, request.(LoginRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Login")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(LoginResponseObject); ok {
		if err := validResponse.VisitLoginResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Register operation middleware
func (sh *strictHandler) Register(w http.ResponseWriter, r *http.Request, params RegisterParams) {
	var request RegisterRequestObject

	request.Params = params

	var body RegisterJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Register(ctx, request.(RegisterRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Register")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(RegisterResponseObject); ok {
		if err := validResponse.VisitRegisterResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// User operation middleware
func (sh *strictHandler) User(w http.ResponseWriter, r *http.Request, userIDParam UserIDParam) {
	var request UserRequestObject

	request.UserIDParam = userIDParam

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.User(ctx, request.(UserRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "User")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(UserResponseObject); ok {
		if err := validResponse.VisitUserResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// UserAvatar operation middleware
func (sh *strictHandler) UserAvatar(w http.ResponseWriter, r *http.Request, userIDParam UserIDParam) {
	var request UserAvatarRequestObject

	request.UserIDParam = userIDParam

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.UserAvatar(ctx, request.(UserAvatarRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "UserAvatar")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(UserAvatarResponseObject); ok {
		if err := validResponse.VisitUserAvatarResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
