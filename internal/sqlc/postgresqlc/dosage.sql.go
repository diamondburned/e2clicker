// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: dosage.sql

package postgresqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	userservice "libdb.so/e2clicker/services/user"
)

const deleteDosageSchedule = `-- name: DeleteDosageSchedule :exec
DELETE FROM dosage_schedule
WHERE user_secret = $1
`

func (q *Queries) DeleteDosageSchedule(ctx context.Context, userSecret userservice.Secret) error {
	_, err := q.db.Exec(ctx, deleteDosageSchedule, userSecret)
	return err
}

const dosageSchedule = `-- name: DosageSchedule :one
/*
 * Dosage and dosage-related
 */
SELECT user_secret, delivery_method, dose, interval, concurrence
FROM dosage_schedule
WHERE user_secret = $1
`

func (q *Queries) DosageSchedule(ctx context.Context, userSecret userservice.Secret) (DosageSchedule, error) {
	row := q.db.QueryRow(ctx, dosageSchedule, userSecret)
	var i DosageSchedule
	err := row.Scan(
		&i.UserSecret,
		&i.DeliveryMethod,
		&i.Dose,
		&i.Interval,
		&i.Concurrence,
	)
	return i, err
}

const doseHistory = `-- name: DoseHistory :many
SELECT dose_id, user_secret, delivery_method, dose, taken_at, taken_off_at
FROM dosage_history
WHERE user_secret = $1
  AND taken_at >= $2
  AND taken_at < $3
  -- order latest last
ORDER BY taken_at ASC
`

type DoseHistoryParams struct {
	UserSecret userservice.Secret
	Start      pgtype.Timestamptz
	End        pgtype.Timestamptz
}

func (q *Queries) DoseHistory(ctx context.Context, arg DoseHistoryParams) ([]DosageHistory, error) {
	rows, err := q.db.Query(ctx, doseHistory, arg.UserSecret, arg.Start, arg.End)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DosageHistory
	for rows.Next() {
		var i DosageHistory
		if err := rows.Scan(
			&i.DoseID,
			&i.UserSecret,
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

const editDose = `-- name: EditDose :execrows
UPDATE
  dosage_history
SET delivery_method = $3, dose = $4, taken_at = $5, taken_off_at = $6
WHERE user_secret = $2
  AND dose_id = $1
RETURNING dose_id, user_secret, delivery_method, dose, taken_at, taken_off_at
`

type EditDoseParams struct {
	DoseID         int64
	UserSecret     userservice.Secret
	DeliveryMethod pgtype.Text
	Dose           float32
	TakenAt        pgtype.Timestamptz
	TakenOffAt     pgtype.Timestamptz
}

func (q *Queries) EditDose(ctx context.Context, arg EditDoseParams) (int64, error) {
	result, err := q.db.Exec(ctx, editDose,
		arg.DoseID,
		arg.UserSecret,
		arg.DeliveryMethod,
		arg.Dose,
		arg.TakenAt,
		arg.TakenOffAt,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const forgetDoses = `-- name: ForgetDoses :execrows
DELETE FROM dosage_history
WHERE user_secret = $1
  AND dose_id = ANY ($2::bigint[])
`

type ForgetDosesParams struct {
	UserSecret userservice.Secret
	DoseIDs    []int64
}

func (q *Queries) ForgetDoses(ctx context.Context, arg ForgetDosesParams) (int64, error) {
	result, err := q.db.Exec(ctx, forgetDoses, arg.UserSecret, arg.DoseIDs)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const recordDose = `-- name: RecordDose :one
INSERT INTO dosage_history (user_secret, taken_at, delivery_method, dose) (
  SELECT $1::usersecret, $2::timestamptz, delivery_method, dose
  FROM dosage_schedule
  WHERE dosage_schedule.user_secret = $1::usersecret)
RETURNING dosage_history.dose_id, dosage_history.user_secret, dosage_history.delivery_method, dosage_history.dose, dosage_history.taken_at, dosage_history.taken_off_at
`

type RecordDoseParams struct {
	UserSecret userservice.Secret
	TakenAt    pgtype.Timestamptz
}

func (q *Queries) RecordDose(ctx context.Context, arg RecordDoseParams) (DosageHistory, error) {
	row := q.db.QueryRow(ctx, recordDose, arg.UserSecret, arg.TakenAt)
	var i DosageHistory
	err := row.Scan(
		&i.DoseID,
		&i.UserSecret,
		&i.DeliveryMethod,
		&i.Dose,
		&i.TakenAt,
		&i.TakenOffAt,
	)
	return i, err
}

const setDosageSchedule = `-- name: SetDosageSchedule :exec
INSERT INTO dosage_schedule (user_secret, delivery_method, dose, interval, concurrence)
  VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_secret)
  DO UPDATE SET
    delivery_method = $2, dose = $3, interval = $4, concurrence = $5
`

type SetDosageScheduleParams struct {
	UserSecret     userservice.Secret
	DeliveryMethod pgtype.Text
	Dose           float32
	Interval       pgtype.Interval
	Concurrence    pgtype.Int2
}

func (q *Queries) SetDosageSchedule(ctx context.Context, arg SetDosageScheduleParams) error {
	_, err := q.db.Exec(ctx, setDosageSchedule,
		arg.UserSecret,
		arg.DeliveryMethod,
		arg.Dose,
		arg.Interval,
		arg.Concurrence,
	)
	return err
}
