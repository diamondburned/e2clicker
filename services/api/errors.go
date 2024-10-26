package api

import (
	"context"

	"libdb.so/e2clicker/internal/publicerrors"
	"libdb.so/e2clicker/services/api/openapi"
)

type errorResponse = struct {
	Body       openapi.Error
	StatusCode int
}

func convertError[T ~errorResponse](ctx context.Context, err error) T {
	return convertErrorWithMessage[T](ctx, err, "")
}

func convertErrorWithMessage[T ~errorResponse](ctx context.Context, err error, hiddenMessage string) T {
	marshaled := publicerrors.MarshalError(ctx, err, hiddenMessage)
	oapiValue := openapi.Error{Message: marshaled.Message}
	if marshaled.Details != nil {
		oapiValue.Details = &marshaled.Details
	}
	if marshaled.Internal {
		oapiValue.Internal = &marshaled.Internal
		oapiValue.InternalCode = &marshaled.InternalCode
	}

	statusCode := 400
	if marshaled.Internal {
		statusCode = 500
	}

	return T{
		Body:       oapiValue,
		StatusCode: statusCode,
	}
}
