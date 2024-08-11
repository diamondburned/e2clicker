// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package postgresqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	xid "github.com/rs/xid"
	assetservice "libdb.so/hrtclicker/v2/services/asset"
	notificationservice "libdb.so/hrtclicker/v2/services/notification"
	userservice "libdb.so/hrtclicker/v2/services/user"
)

type Compression string

const (
	CompressionGzip   Compression = "gzip"
	CompressionZstd   Compression = "zstd"
	CompressionBrotli Compression = "brotli"
)

func (e *Compression) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Compression(s)
	case string:
		*e = Compression(s)
	default:
		return fmt.Errorf("unsupported scan type for Compression: %T", src)
	}
	return nil
}

type NullCompression struct {
	Compression Compression
	Valid       bool // Valid is true if Compression is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCompression) Scan(value interface{}) error {
	if value == nil {
		ns.Compression, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Compression.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCompression) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Compression), nil
}

type DeliveryMethod struct {
	ID    string
	Units string
	Name  string
}

type DosageHistory struct {
	DoseID         int64
	LastDose       pgtype.Int8
	UserID         pgtype.Uint32
	DeliveryMethod pgtype.Text
	Dose           pgtype.Numeric
	TakenAt        pgtype.Timestamptz
	TakenOffAt     pgtype.Timestamptz
}

type DosageSchedule struct {
	UserID         xid.ID
	DeliveryMethod pgtype.Text
	Dose           pgtype.Numeric
	Interval       pgtype.Interval
	Concurrence    pgtype.Numeric
}

type Meta struct {
	X bool
	V int16
}

type NotificationHistory struct {
	NotificationID int64
	UserID         pgtype.Uint32
	DosageID       pgtype.Int8
	SentAt         pgtype.Timestamptz
	ErrorReason    pgtype.Text
}

type User struct {
	UserID              xid.ID
	Email               string
	Passhash            []byte
	Name                string
	Locale              userservice.Locale
	RegisteredAt        pgtype.Timestamp
	NotificationService *notificationservice.NotificationConfigJSON
	CustomNotification  *notificationservice.Notification
}

type UserAvatar struct {
	UserID      xid.ID
	MimeType    string
	Compression assetservice.Compression
	AvatarImage []byte
}

type UserSession struct {
	SessionID int64
	UserID    xid.ID
	Token     []byte
	CreatedAt pgtype.Timestamp
	LastUsed  pgtype.Timestamp
	UserAgent pgtype.Text
}