package user

import "github.com/rs/xid"

// UserID is a unique identifier for a user.
type UserID xid.ID

// User is a user in the system.
type User struct {
	ID        UserID `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	HasAvatar bool   `json:"has_avatar"`
}
