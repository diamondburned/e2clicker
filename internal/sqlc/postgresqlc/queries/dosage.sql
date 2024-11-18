/*
 * Dosage and dosage-related
 */
-- name: DosageSchedule :one
SELECT *
FROM dosage_schedule
WHERE user_secret = $1;

-- name: SetDosageSchedule :exec
INSERT INTO dosage_schedule (user_secret, delivery_method, dose, interval, concurrence)
  VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_secret)
  DO UPDATE SET
    delivery_method = $2, dose = $3, interval = $4, concurrence = $5;

-- name: DeleteDosageSchedule :exec
DELETE FROM dosage_schedule
WHERE user_secret = $1;

-- name: RecordDose :one
INSERT INTO dosage_history (user_secret, taken_at, delivery_method, dose) (
  SELECT @user_secret::usersecret, @taken_at::timestamptz, delivery_method, dose
  FROM dosage_schedule
  WHERE dosage_schedule.user_secret = @user_secret::usersecret)
RETURNING dosage_history.*;

-- name: EditDose :execrows
UPDATE
  dosage_history
SET delivery_method = $3, dose = $4, taken_at = $5, taken_off_at = $6
WHERE user_secret = $2
  AND dose_id = $1
RETURNING *;

-- name: ForgetDoses :execrows
DELETE FROM dosage_history
WHERE user_secret = $1
  AND dose_id = ANY (@dose_ids::bigint[]);

-- name: DoseHistory :many
SELECT *
FROM dosage_history
WHERE user_secret = $1
  AND taken_at >= sqlc.arg('start')
  AND taken_at < sqlc.arg('end')
  -- order latest last
ORDER BY taken_at ASC;
