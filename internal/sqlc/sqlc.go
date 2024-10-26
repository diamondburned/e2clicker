package sqlc

import "libdb.so/e2clicker/services/user"

//go:generate sqlc generate

// UserID is a wrapper around user.UserID to add SQL type methods.
type UserID struct {
	user.UserID
}

func (id *UserID) ScanBytes(v []byte) error {
	u, err := user.ParseRawUserID(v)
	if err != nil {
		return err
	}
	id.UserID = u
	return nil
}

func (id UserID) BytesValue() ([]byte, error) {
	return id.UserID.Bytes(), nil
}
