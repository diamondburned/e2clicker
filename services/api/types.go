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

func maybeNil[T any](value T, valid bool) *T {
	if valid {
		return &value
	}
	return nil
}

func convertUser(u user.User) openapi.User {
	return openapi.User(u)
}

func convertSession(s user.Session) openapi.Session {
	return openapi.Session{
		ID:        s.ID,
		CreatedAt: s.CreatedAt,
		LastUsed:  s.LastUsed,
		ExpiresAt: s.ExpiresAt,
	}
}

func convertList[T any, U any](list []T, convert func(T) U) []U {
	result := make([]U, len(list))
	for i, item := range list {
		result[i] = convert(item)
	}
	return result
}
