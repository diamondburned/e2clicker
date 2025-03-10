// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package postgresqlc

import (
	notificationservice "e2clicker.app/services/notification"
	userservice "e2clicker.app/services/user"
	"github.com/jackc/pgx/v5/pgtype"
)

type DeliveryMethod struct {
	ID          string
	Units       string
	Name        string
	Description string
}

type DosageHistory struct {
	UserSecret     userservice.Secret
	DeliveryMethod pgtype.Text
	Dose           float32
	TakenAt        pgtype.Timestamptz
	TakenOffAt     pgtype.Timestamptz
	Comment        pgtype.Text
}

type DosageSchedule struct {
	UserSecret     userservice.Secret
	DeliveryMethod pgtype.Text
	Dose           float32
	Interval       pgtype.Interval
	Concurrence    pgtype.Int2
}

type Meta struct {
	X bool
	V int16
}

type NotificationHistory struct {
	NotificationID     pgtype.UUID
	UserSecret         userservice.Secret
	SupposedEntityTime pgtype.Timestamptz
	SentAt             pgtype.Timestamptz
	ErrorReason        pgtype.Text
	Errored            pgtype.Bool
}

type User struct {
	Secret                  userservice.Secret
	Name                    string
	Locale                  userservice.Locale
	RegisteredAt            pgtype.Timestamp
	NotificationPreferences notificationservice.UserPreferences
}

type UserSession struct {
	ID         int64
	UserSecret userservice.Secret
	Token      []byte
	CreatedAt  pgtype.Timestamp
	LastUsed   pgtype.Timestamp
	UserAgent  pgtype.Text
}
