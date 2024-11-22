// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version (devel) DO NOT EDIT.
package openapi

// Error defines model for Error.
type Error struct {
	// Details Additional details about the error. Ignored if [errors] is used.
	Details *interface{} `json:"details,omitempty"`

	// Errors An array of errors that caused this error. If this is populated, then [details] is omitted.
	Errors []Error `json:"errors,omitempty"`

	// Internal Whether the error is internal
	Internal *bool `json:"internal,omitempty"`

	// InternalCode An internal code for the error (useless for clients)
	InternalCode *string `json:"internalCode,omitempty"`

	// Message A message describing the error
	Message string `json:"message"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse = Error

// RateLimitedResponse defines model for RateLimitedResponse.
type RateLimitedResponse = Error
