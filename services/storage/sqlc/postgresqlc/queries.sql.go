// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package postgresqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	xid "github.com/rs/xid"
	notificationservice "libdb.so/hrtclicker/v2/services/notification"
	userservice "libdb.so/hrtclicker/v2/services/user"
)

const addUserDosage = `-- name: AddUserDosage :exec
INSERT INTO dosage_history (user_id, last_dose, taken_at, delivery_method, dose) (
  SELECT $1, $2, now(), delivery_method, dose
  FROM dosage_schedule
  WHERE user_id = $1)
`

type AddUserDosageParams struct {
	UserID   pgtype.Uint32
	LastDose pgtype.Int8
}

func (q *Queries) AddUserDosage(ctx context.Context, arg AddUserDosageParams) error {
	_, err := q.db.Exec(ctx, addUserDosage, arg.UserID, arg.LastDose)
	return err
}

const createUser = `-- name: CreateUser :exec
/*
 * User
 */
INSERT INTO users (user_id, email, passhash, name)
  VALUES ($1, $2, $3, $4)
`

type CreateUserParams struct {
	UserID   xid.ID
	Email    string
	Passhash []byte
	Name     string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.UserID,
		arg.Email,
		arg.Passhash,
		arg.Name,
	)
	return err
}

const deliveryMethod = `-- name: DeliveryMethod :one
SELECT id, units, name
FROM delivery_methods
WHERE name = $1
`

func (q *Queries) DeliveryMethod(ctx context.Context, name string) (DeliveryMethod, error) {
	row := q.db.QueryRow(ctx, deliveryMethod, name)
	var i DeliveryMethod
	err := row.Scan(&i.ID, &i.Units, &i.Name)
	return i, err
}

const deliveryMethods = `-- name: DeliveryMethods :many
/*
 * Delivery Method
 */
SELECT id, units, name
FROM delivery_methods
`

func (q *Queries) DeliveryMethods(ctx context.Context) ([]DeliveryMethod, error) {
	rows, err := q.db.Query(ctx, deliveryMethods)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DeliveryMethod
	for rows.Next() {
		var i DeliveryMethod
		if err := rows.Scan(&i.ID, &i.Units, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const registerSession = `-- name: RegisterSession :exec
/*
 * User Session
 */
INSERT INTO user_sessions (user_id, token, created_at, last_used, user_agent)
  VALUES ($1, $2, now(), now(), $3)
`

type RegisterSessionParams struct {
	UserID    xid.ID
	Token     []byte
	UserAgent pgtype.Text
}

func (q *Queries) RegisterSession(ctx context.Context, arg RegisterSessionParams) error {
	_, err := q.db.Exec(ctx, registerSession, arg.UserID, arg.Token, arg.UserAgent)
	return err
}

const setUserAvatar = `-- name: SetUserAvatar :exec
INSERT INTO user_avatars (user_id, avatar_image)
  VALUES ($1, $2)
ON CONFLICT (user_id)
  DO UPDATE SET
    avatar_image = $2
`

type SetUserAvatarParams struct {
	UserID      xid.ID
	AvatarImage []byte
}

func (q *Queries) SetUserAvatar(ctx context.Context, arg SetUserAvatarParams) error {
	_, err := q.db.Exec(ctx, setUserAvatar, arg.UserID, arg.AvatarImage)
	return err
}

const setUserCustomNotification = `-- name: SetUserCustomNotification :exec
UPDATE
  users
SET custom_notification = $2
WHERE user_id = $1
`

type SetUserCustomNotificationParams struct {
	UserID             xid.ID
	CustomNotification *notificationservice.Notification
}

func (q *Queries) SetUserCustomNotification(ctx context.Context, arg SetUserCustomNotificationParams) error {
	_, err := q.db.Exec(ctx, setUserCustomNotification, arg.UserID, arg.CustomNotification)
	return err
}

const setUserNotificationService = `-- name: SetUserNotificationService :exec
/*
 * User Notifications
 */
UPDATE
  users
SET notification_service = $2
WHERE user_id = $1
`

type SetUserNotificationServiceParams struct {
	UserID              xid.ID
	NotificationService *notificationservice.NotificationConfigJSON
}

func (q *Queries) SetUserNotificationService(ctx context.Context, arg SetUserNotificationServiceParams) error {
	_, err := q.db.Exec(ctx, setUserNotificationService, arg.UserID, arg.NotificationService)
	return err
}

const updateUserDosageSchedule = `-- name: UpdateUserDosageSchedule :exec
INSERT INTO dosage_schedule (user_id, delivery_method, dose, interval, concurrence)
  VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id)
  DO UPDATE SET
    delivery_method = $2, dose = $3, interval = $4, concurrence = $5
`

type UpdateUserDosageScheduleParams struct {
	UserID         xid.ID
	DeliveryMethod pgtype.Text
	Dose           pgtype.Numeric
	Interval       pgtype.Interval
	Concurrence    pgtype.Numeric
}

func (q *Queries) UpdateUserDosageSchedule(ctx context.Context, arg UpdateUserDosageScheduleParams) error {
	_, err := q.db.Exec(ctx, updateUserDosageSchedule,
		arg.UserID,
		arg.DeliveryMethod,
		arg.Dose,
		arg.Interval,
		arg.Concurrence,
	)
	return err
}

const updateUserEmailPassword = `-- name: UpdateUserEmailPassword :exec
UPDATE
  users
SET email = $2, passhash = $3
WHERE user_id = $1
`

type UpdateUserEmailPasswordParams struct {
	UserID   xid.ID
	Email    string
	Passhash []byte
}

func (q *Queries) UpdateUserEmailPassword(ctx context.Context, arg UpdateUserEmailPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserEmailPassword, arg.UserID, arg.Email, arg.Passhash)
	return err
}

const updateUserLocale = `-- name: UpdateUserLocale :exec
UPDATE
  users
SET locale = $2
WHERE user_id = $1
`

type UpdateUserLocaleParams struct {
	UserID xid.ID
	Locale userservice.Locale
}

func (q *Queries) UpdateUserLocale(ctx context.Context, arg UpdateUserLocaleParams) error {
	_, err := q.db.Exec(ctx, updateUserLocale, arg.UserID, arg.Locale)
	return err
}

const updateUserName = `-- name: UpdateUserName :exec
UPDATE
  users
SET name = $2
WHERE user_id = $1
`

type UpdateUserNameParams struct {
	UserID xid.ID
	Name   string
}

func (q *Queries) UpdateUserName(ctx context.Context, arg UpdateUserNameParams) error {
	_, err := q.db.Exec(ctx, updateUserName, arg.UserID, arg.Name)
	return err
}

const user = `-- name: User :one
SELECT user_id, email, name, EXISTS (
    SELECT user_id
    FROM user_avatars
    WHERE user_id = users.user_id) AS has_avatar
FROM users
WHERE users.user_id = $1
`

type UserRow struct {
	UserID    xid.ID
	Email     string
	Name      string
	HasAvatar bool
}

func (q *Queries) User(ctx context.Context, userID xid.ID) (UserRow, error) {
	row := q.db.QueryRow(ctx, user, userID)
	var i UserRow
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Name,
		&i.HasAvatar,
	)
	return i, err
}

const userAvatar = `-- name: UserAvatar :one
/*
 * User Avatar
 */
SELECT avatar_image
FROM user_avatars
WHERE user_id = $1
`

func (q *Queries) UserAvatar(ctx context.Context, userID xid.ID) ([]byte, error) {
	row := q.db.QueryRow(ctx, userAvatar, userID)
	var avatar_image []byte
	err := row.Scan(&avatar_image)
	return avatar_image, err
}

const userDosageHistory = `-- name: UserDosageHistory :many
SELECT dose_id, last_dose, user_id, delivery_method, dose, taken_at, taken_off_at
FROM dosage_history
WHERE user_id = $1
  AND taken_at >= $2
ORDER BY dose_id DESC
LIMIT $3
`

type UserDosageHistoryParams struct {
	UserID  pgtype.Uint32
	TakenAt pgtype.Timestamptz
	Limit   int32
}

func (q *Queries) UserDosageHistory(ctx context.Context, arg UserDosageHistoryParams) ([]DosageHistory, error) {
	rows, err := q.db.Query(ctx, userDosageHistory, arg.UserID, arg.TakenAt, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DosageHistory
	for rows.Next() {
		var i DosageHistory
		if err := rows.Scan(
			&i.DoseID,
			&i.LastDose,
			&i.UserID,
			&i.DeliveryMethod,
			&i.Dose,
			&i.TakenAt,
			&i.TakenOffAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const userDosageSchedule = `-- name: UserDosageSchedule :exec
/*
 * Dosage and dosage-related
 */
SELECT user_id, delivery_method, dose, interval, concurrence
FROM dosage_schedule
WHERE user_id = $1
`

func (q *Queries) UserDosageSchedule(ctx context.Context, userID xid.ID) error {
	_, err := q.db.Exec(ctx, userDosageSchedule, userID)
	return err
}

const userPasswordHashFromEmail = `-- name: UserPasswordHashFromEmail :one
SELECT user_id, passhash
FROM users
WHERE email = $1
`

type UserPasswordHashFromEmailRow struct {
	UserID   xid.ID
	Passhash []byte
}

func (q *Queries) UserPasswordHashFromEmail(ctx context.Context, email string) (UserPasswordHashFromEmailRow, error) {
	row := q.db.QueryRow(ctx, userPasswordHashFromEmail, email)
	var i UserPasswordHashFromEmailRow
	err := row.Scan(&i.UserID, &i.Passhash)
	return i, err
}

const validateSession = `-- name: ValidateSession :one
UPDATE
  user_sessions
SET last_used = now()
FROM users
WHERE user_sessions.user_id = users.user_id
  AND token = $1
  AND last_used > now() - '7 days'::interval
RETURNING users.user_id, users.email, users.name, EXISTS (
    SELECT user_id
    FROM user_avatars
    WHERE user_id = users.user_id) AS has_avatar
`

type ValidateSessionRow struct {
	UserID    xid.ID
	Email     string
	Name      string
	HasAvatar bool
}

func (q *Queries) ValidateSession(ctx context.Context, token []byte) (ValidateSessionRow, error) {
	row := q.db.QueryRow(ctx, validateSession, token)
	var i ValidateSessionRow
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Name,
		&i.HasAvatar,
	)
	return i, err
}

const version = `-- name: Version :one
/*
 * Meta
 */
SELECT v
FROM meta
`

func (q *Queries) Version(ctx context.Context) (int16, error) {
	row := q.db.QueryRow(ctx, version)
	var v int16
	err := row.Scan(&v)
	return v, err
}