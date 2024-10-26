// Package commonsvc defines common services for dependency injection.
package commonsvc

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/samber/do/v2"
)

// HTTPClient provides an HTTP client.
var HTTPClient = do.Lazy(func(do.Injector) (*http.Client, error) {
	return http.DefaultClient, nil
})

// Logger provides a logger.
var Logger = do.Lazy(func(do.Injector) (*slog.Logger, error) {
	return slog.Default(), nil
})

// Context provides a context.
var Context = do.Eager(context.Background())
