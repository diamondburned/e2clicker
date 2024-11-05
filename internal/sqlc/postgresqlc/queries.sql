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
-- name: UserDosageSchedule :exec
SELECT *
FROM dosage_schedule
WHERE user_secret = $1;

-- name: UpdateUserDosageSchedule :exec
INSERT INTO dosage_schedule (user_secret, delivery_method, dose, interval, concurrence)
  VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_secret)
  DO UPDATE SET
    delivery_method = $2, dose = $3, interval = $4, concurrence = $5;

-- name: UserDosageHistory :many
SELECT *
FROM dosage_history
WHERE user_secret = $1
  AND taken_at >= $2
ORDER BY dose_id DESC
LIMIT $3;

-- name: AddUserDosage :exec
INSERT INTO dosage_history (user_secret, last_dose, taken_at, delivery_method, dose) (
  SELECT $1, $2, now(), delivery_method, dose
  FROM dosage_schedule
  WHERE user_secret = $1);
