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
		newOpenAPIHandler,
		newOpenAPIHandlerForImportExport,

		(*Authenticator).AuthenticationFunc, // openapi3filter.AuthenticationFunc
		// (*openAPIHandler).asHandler,                // openapi.ServerInterface
		(*openAPIHandlerForImportExport).asHandler, // openapi.ServerInterface
	),
)
