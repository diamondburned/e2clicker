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

-- name: RecordDose :exec
INSERT INTO dosage_history (user_secret, delivery_method, dose, taken_at, taken_off_at)
  VALUES ($1, $2, $3, $4, $5);

-- name: EditDose :execrows
UPDATE
  dosage_history
SET delivery_method = @delivery_method, dose = @dose, taken_at = @taken_at, taken_off_at = @taken_off_at
WHERE user_secret = @user_secret
  AND taken_at = @old_taken_at;

-- name: ForgetDoses :execrows
DELETE FROM dosage_history
WHERE user_secret = $1
  AND taken_at = ANY (@taken_at::timestamp[]);

-- name: DoseHistory :iter
SELECT *
FROM dosage_history
WHERE user_secret = $1
  AND taken_at >= sqlc.arg('start')
  AND taken_at < sqlc.arg('end')
  -- order latest last
ORDER BY taken_at ASC;
