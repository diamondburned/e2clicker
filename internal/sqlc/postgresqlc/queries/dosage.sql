/*
 * Dosage and dosage-related
 */
-- name: DosageSchedule :one
SELECT *
FROM dosage_schedule
WHERE user_secret = $1;

-- name: SetDosageSchedule :exec
INSERT INTO dosage_schedule (user_secret, delivery_method, dose, interval, concurrence, reminder_recurrence)
  VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (user_secret)
  DO UPDATE SET
    delivery_method = $2, dose = $3, interval = $4, concurrence = $5, reminder_recurrence = $6;

-- name: DeleteDosageSchedule :exec
DELETE FROM dosage_schedule
WHERE user_secret = $1;

-- name: RecordDose :one
INSERT INTO dosage_history (user_secret, id, delivery_method, dose, taken_at, taken_off_at, comment)
  VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: EditDose :execrows
UPDATE
  dosage_history
SET delivery_method = @delivery_method, dose = @dose, taken_at = @taken_at, taken_off_at = @taken_off_at, comment = @comment
WHERE user_secret = @user_secret
  AND id = @id;

-- name: ForgetDoses :execrows
DELETE FROM dosage_history
WHERE user_secret = $1
  AND id = ANY (@ids::uuid[]);

-- name: DoseHistory :many
SELECT *
FROM dosage_history
WHERE user_secret = $1
  AND taken_at >= sqlc.arg('start')
  AND taken_at < sqlc.arg('end')
  -- order latest last
ORDER BY taken_at ASC, id ASC;

-- name: UpcomingDosageReminders :many
SELECT
  users.secret AS user_secret,
  users.name AS user_name,
  sqlc.embed(dosage_schedule),
  sqlc.embed(dosage_history),
  last_notification.supposed_entity_time AS last_reminded_dose,
  last_notification.sent_at AS last_reminded_at
FROM users
  INNER JOIN dosage_schedule ON users.secret = dosage_schedule.user_secret
  INNER JOIN LATERAL (
    SELECT *
    FROM dosage_history dh
    WHERE dh.user_secret = users.secret
    ORDER BY dh.taken_at DESC, dh.id DESC
    LIMIT 1
  ) dosage_history ON TRUE
  LEFT JOIN LATERAL (
    SELECT nh.supposed_entity_time, nh.sent_at
    FROM notification_history nh
    WHERE nh.user_secret = users.secret
    ORDER BY nh.supposed_entity_time DESC NULLS LAST, nh.sent_at DESC
    LIMIT 1
  ) last_notification ON TRUE;

-- name: RecordRemindedDoseAttempt :exec
INSERT INTO notification_history (user_secret, sent_at, supposed_entity_time, error_reason)
  VALUES ($1, $2, $3, $4);
