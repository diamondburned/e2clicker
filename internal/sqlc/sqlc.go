package sqlc

import (
	"github.com/rs/xid"
)

//go:generate sqlc generate

// XID is a wrapper around xid.ID to add SQL type methods.
type XID xid.ID

func (id *XID) ScanBytes(b []byte) error {
	v, err := xid.FromBytes(b)
	if err != nil {
		return err
	}
	*id = XID(v)
	return nil
}

func (id XID) BytesValue() ([]byte, error) {
	return xid.ID(id).Bytes(), nil
}
