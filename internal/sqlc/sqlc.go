package sqlc

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/xid"
)

//go:generate sqlc generate

// XID is a wrapper around xid.ID to add SQL type methods.
type XID xid.ID

var (
	_ pgtype.TextScanner = (*XID)(nil)
	_ pgtype.TextValuer  = (*XID)(nil)
)

func (id *XID) ScanText(t pgtype.Text) error {
	if !t.Valid {
		*id = XID{}
		return nil
	}

	v, err := xid.FromString(t.String)
	if err != nil {
		return fmt.Errorf("can't scan as text: %w", err)
	}
	*id = XID(v)
	return nil
}

func (id XID) TextValue() (pgtype.Text, error) {
	return pgtype.Text{
		String: xid.ID(id).String(),
		Valid:  !xid.ID(id).IsZero(),
	}, nil
}
