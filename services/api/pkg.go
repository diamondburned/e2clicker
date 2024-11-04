package api

import (
	"log/slog"

	"go.uber.org/fx"
)

var Module = fx.Module("api",
	fx.Decorate(func(slog *slog.Logger) *slog.Logger {
		return slog.With("module", "api")
	}),
	fx.Provide(
		NewServer,
		NewAuthenticator,
		NewOpenAPIHandler,

		(*Authenticator).AuthenticationFunc, // openapi3filter.AuthenticationFunc
		(*OpenAPIHandler).asStrictHandler,   // openapi.StrictServerInterface
	),
)
