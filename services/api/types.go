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

func convertUser(u user.User) openapi.User {
	return openapi.User{
		ID:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		Locale: u.Locale,
	}
}
