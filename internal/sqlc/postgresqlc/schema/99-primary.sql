-- Version 0:
-- This can be run for as many time as is needed.
CREATE TABLE IF NOT EXISTS meta (
  x bool PRIMARY KEY DEFAULT TRUE CHECK (x), -- force only 1 row
  v smallint NOT NULL
);

INSERT INTO meta (v)
  VALUES (1)
ON CONFLICT
  DO NOTHING;

-- NEW VERSION
UPDATE
  meta
SET v = 2;

CREATE DOMAIN usersecret AS text COLLATE "C";

CREATE DOMAIN notificationpreferences AS jsonb;

CREATE DOMAIN locale AS text;

CREATE TABLE users (
  -- The user's secret (a random string).
  secret usersecret PRIMARY KEY,
  -- The user's display name.
  name text NOT NULL,
  -- The user's locale.
  locale locale NOT NULL DEFAULT 'en-US',
  -- The time the user was created.
  registered_at timestamp NOT NULL DEFAULT now(),
  -- The [notification.UserPreferences] type in the Go codebase.
  notification_preferences notificationpreferences NOT NULL DEFAULT cast('{}' AS jsonb)
);

CREATE INDEX users_secret ON users USING HASH (secret);

CREATE TABLE user_sessions (
  -- The session ID. This should never be used to log in, but it can be used
  -- to revoke a session.
  id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  -- The user's secret.
  user_secret usersecret NOT NULL REFERENCES users (secret) ON DELETE CASCADE,
  -- The session token.
  token bytea UNIQUE NOT NULL,
  -- The time the session was created.
  created_at timestamp NOT NULL DEFAULT now(),
  -- The time the session was last used.
  last_used timestamp NOT NULL DEFAULT now(),
  -- The user agent string, if any. Collected for management purposes.
  user_agent text
);

CREATE INDEX user_sessions_user_secret ON user_sessions USING HASH (user_secret);

CREATE INDEX user_sessions_token ON user_sessions USING HASH (token);

CREATE TABLE dosage_schedule (
  user_secret usersecret PRIMARY KEY REFERENCES users (secret) ON DELETE CASCADE,
  -- The delivery method of the medication.
  delivery_method text REFERENCES delivery_methods (id) ON DELETE CASCADE,
  -- The dose of the medication.
  -- The units are determined by the delivery method.
  dose real NOT NULL,
  -- The interval between doses.
  interval interval NOT NULL CHECK (interval > '0 minutes'::interval),
  -- How many patches are on at a time. Only relevant (non-null) for patches.
  concurrence smallint
);

CREATE TABLE dosage_history (
  -- The user that took the dose.
  user_secret usersecret NOT NULL REFERENCES users (secret) ON DELETE CASCADE,
  -- The delivery method of the medication.
  delivery_method text REFERENCES delivery_methods (id) ON DELETE CASCADE,
  -- The dose of the medication.
  -- This is usually copied from the schedule, but can be overridden.
  dose real NOT NULL,
  -- The time the dose was taken.
  taken_at timestamptz NOT NULL,
  -- The time the patch was taken off. This is only applicable to patches.
  taken_off_at timestamptz,
  -- The comment for the dose.
  comment text,
  -- Ensure that the user can't take multiple doses at the same time.
  -- Realistically, they can, but our system does not allow this anyway unless
  -- you're trying to import in bulk.
  UNIQUE (user_secret, taken_at)
);

CREATE INDEX dosage_history_user_secret ON dosage_history USING HASH (user_secret);

CREATE INDEX dosage_history_taken_at ON dosage_history USING BTREE (taken_at);

CREATE TABLE notification_history (
  notification_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  -- The user that the notification was for.
  user_secret usersecret NOT NULL REFERENCES users (secret) ON DELETE CASCADE,
  -- The timestamp of the entity that the notification is about.
  -- For dosage notifications, this is the time the dose is supposed to be taken.
  supposed_entity_time timestamptz,
  -- The time the notification was sent. This indicates an attempt.
  sent_at timestamptz NOT NULL,
  -- The error if the notification failed to send.
  error_reason text,
  -- True if the notification errored.
  errored boolean GENERATED ALWAYS AS (error_reason IS NOT NULL) STORED
);
