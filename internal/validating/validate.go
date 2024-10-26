package validating

import "context"

// Validator is the interface that wraps the Validate method.
type Validator interface {
	Validate() error
}

// ContextValidator is a variant of the Validator interface that can validate
// with a context.
type ContextValidator interface {
	ValidateContext(ctx context.Context) error
}

// Validate validates the given value.
// If v does not implement either [Validator] or [ContextValidator], this
// function does nothing.
func Validate(ctx context.Context, v any) error {
	if validator, ok := v.(Validator); ok {
		return validator.Validate()
	}
	if validator, ok := v.(ContextValidator); ok {
		return validator.ValidateContext(ctx)
	}
	return nil
}
