package user

import "testing"

//go:generate moq -out user_mock_test.go -stub . UserStorage UserAvatarStorage UserSessionStorage

type mockUserService struct {
	UserService
	users    *UserStorageMock
	avatars  *UserAvatarStorageMock
	sessions *UserSessionStorageMock
}

func newMockUserService(*testing.T) *mockUserService {
	s := &mockUserService{
		users:    &UserStorageMock{},
		avatars:  &UserAvatarStorageMock{},
		sessions: &UserSessionStorageMock{},
	}
	s.UserService = UserService{
		s.users,
		s.avatars,
		s.sessions,
	}
	return s
}
