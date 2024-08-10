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
-- name: CreateUser :exec
INSERT INTO users (user_id, email, passhash, name)
  VALUES ($1, $2, $3, $4);

-- name: UserByEmail :one
SELECT user_id, email, name, EXISTS (
    SELECT user_id
    FROM user_avatars
    WHERE user_id = users.user_id) AS has_avatar
FROM users
WHERE email = $1;

-- name: UserByID :one
SELECT user_id, email, name, EXISTS (
    SELECT user_id
    FROM user_avatars
    WHERE user_id = users.user_id) AS has_avatar
FROM users
WHERE users.user_id = $1;

-- name: UserConfigureNotifications :exec
UPDATE
  users
SET notification_service = $2
WHERE user_id = $1;

-- name: UserSetNotificationMessage :exec
UPDATE
  users
SET notification_message = $2
WHERE user_id = $1;

-- name: UserUpdateEmailPassword :exec
UPDATE
  users
SET email = $2, passhash = $3
WHERE user_id = $1;

-- name: UserUpdateName :exec
UPDATE
  users
SET name = $2
WHERE user_id = $1;

-- name: UserAvatar :one
SELECT avatar_image
FROM user_avatars
WHERE user_id = $1;

-- name: UserSetAvatar :exec
INSERT INTO user_avatars (user_id, avatar_image)
  VALUES ($1, $2)
ON CONFLICT (user_id)
  DO UPDATE SET
    avatar_image = $2;

-- name: RegisterSession :exec
INSERT INTO user_sessions (user_id, token, created_at, last_used, user_agent)
  VALUES ($1, $2, now(), now(), $3);

-- name: UserValidateSession :one
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
    WHERE user_id = users.user_id) AS has_avatar;


/*
 * Dosage and dosage-related
 */
-- name: UserDosageSchedule :exec
SELECT *
FROM dosage_schedule
WHERE user_id = $1;

-- name: UserUpdateDosageSchedule :exec
INSERT INTO dosage_schedule (user_id, delivery_method, dose, interval, concurrence)
  VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id)
  DO UPDATE SET
    delivery_method = $2, dose = $3, interval = $4, concurrence = $5;

-- name: UserDosageHistory :many
SELECT *
FROM dosage_history
WHERE user_id = $1
  AND taken_at >= $2
ORDER BY dose_id DESC
LIMIT $3;

-- name: UserAddDosage :exec
INSERT INTO dosage_history (user_id, last_dose, taken_at, delivery_method, dose) (
  SELECT $1, $2, now(), delivery_method, dose
  FROM dosage_schedule
  WHERE user_id = $1);
