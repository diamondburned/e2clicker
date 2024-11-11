/*
 * Meta
 */
-- name: Version :one
SELECT v
FROM meta;


/*
 * Delivery Method
 */
-- name: DeliveryMethods :many
SELECT *
FROM delivery_methods;

-- name: DeliveryMethod :one
SELECT *
FROM delivery_methods
WHERE name = $1;


/*
 * User
 */
-- name: CreateUser :one
INSERT INTO users (secret, name)
  VALUES ($1, $2)
RETURNING *;

-- name: User :one
SELECT *
FROM users_with_avatar
WHERE secret = $1;

-- name: UpdateUserName :exec
UPDATE
  users
SET name = $2
WHERE secret = $1;

-- name: UpdateUserLocale :exec
UPDATE
  users
SET locale = $2
WHERE secret = $1;


/*
 * User Notifications
 */
-- name: UserNotificationPreferences :one
SELECT notification_preferences
FROM users
WHERE secret = $1;

-- name: SetUserNotificationPreferences :exec
UPDATE
  users
SET notification_preferences = $2
WHERE secret = $1;


/*
 * User Avatar
 */
-- name: UserAvatar :one
SELECT avatar_image, mime_type
FROM user_avatars
WHERE user_secret = $1;

-- name: SetUserAvatar :exec
INSERT INTO user_avatars (user_secret, avatar_image, mime_type)
  VALUES ($1, $2, $3)
ON CONFLICT (user_secret)
  DO UPDATE SET
    avatar_image = $2, mime_type = $3;


/*                                                                                 
 * User Session                                                                    
 */
-- name: RegisterSession :exec
INSERT INTO user_sessions (user_secret, token, created_at, last_used, user_agent)
  VALUES ($1, $2, now(), now(), $3);

-- name: ValidateSession :one
UPDATE
  user_sessions
SET last_used = now()
WHERE token = $1
RETURNING *;

-- name: ListSessions :many
SELECT *
FROM user_sessions
WHERE user_secret = $1;

-- name: DeleteSession :exec
DELETE FROM user_sessions
WHERE user_secret = $1
  AND id = $2;


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
  SELECT @user_secret::xid_, @taken_at::timestamptz, delivery_method, dose
  FROM dosage_schedule
  WHERE dosage_schedule.user_secret = @user_secret::xid_)
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
