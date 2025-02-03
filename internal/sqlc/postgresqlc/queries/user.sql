/*
 * User
 */
-- name: CreateUser :one
INSERT INTO users (secret, name)
  VALUES ($1, $2)
RETURNING *;

-- name: User :one
SELECT *
FROM users
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
SET notification_preferences = $2::jsonb
WHERE secret = $1;


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
