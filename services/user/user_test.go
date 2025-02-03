package user

import "testing"

//go:generate moq -out user_mock_test.go -stub . UserStorage UserSessionStorage

type mockUserService struct {
	UserService
	users    *UserStorageMock
	sessions *UserSessionStorageMock
}

func newMockUserService(*testing.T) *mockUserService {
	s := &mockUserService{
		users:    &UserStorageMock{},
		sessions: &UserSessionStorageMock{},
	}
	s.UserService = UserService{
		s.users,
		s.sessions,
	}
	return s
}
