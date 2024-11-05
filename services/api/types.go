package api

import (
	"libdb.so/e2clicker/services/api/openapi"
	"libdb.so/e2clicker/services/user"
)

func optstr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func optPtr[T any](s *T) T {
	if s != nil {
		return *s
	}
	var v T
	return v
}

func anyPtr[T any](value T) *any {
	v := any(value)
	return &v
}

func convertUser(u user.User) openapi.User {
	return openapi.User(u)
}
