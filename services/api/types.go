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

func anyPtr[T any](value T) *any {
	v := any(value)
	return &v
}

func convertUser(u user.User) openapi.User {
	return openapi.User(u)
}
