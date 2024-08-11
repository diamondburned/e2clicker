-- Language: postgresql
--
CREATE TABLE meta (
  x bool PRIMARY KEY DEFAULT TRUE CHECK (x), -- force only 1 row
  v smallint NOT NULL
);

INSERT INTO meta (v)
  VALUES (1);

-- NEW VERSION
UPDATE
  meta
SET v = 2;

CREATE DOMAIN xid AS char(20);

CREATE DOMAIN assetid AS char(20);

CREATE DOMAIN notification AS jsonb;

CREATE DOMAIN notificationconfig AS jsonb;

CREATE DOMAIN locale AS text;

CREATE TYPE compression AS ENUM (
  'gzip', 'zstd', 'brotli'
);

CREATE TABLE delivery_methods (
  id text PRIMARY KEY, units text NOT NULL, name text NOT NULL
);

INSERT INTO delivery_methods (id, units, name)
  VALUES ('EB im', 'mg', 'Estradiol Benzoate, Intramuscular'),
  ('EV im', 'mg', 'Estradiol Valerate, Intramuscular'),
  ('EEn im', 'mg', 'Estradiol Enanthate, Intramuscular'),
  ('EC im', 'mg', 'Estradiol Cypionate, Intramuscular'),
  ('EUn im', 'mg', 'Estradiol Undecylate, Intramuscular'),
  ('EUn casubq', 'mg', 'Estradiol Undecylate in Castor oil, Subcutaneous'),
  ('patch', 'mcg/day', 'Patch');

CREATE TABLE users (
  -- The user's ID in xid format.
  user_id xid PRIMARY KEY,
  -- The user's email address.
  email text NOT NULL UNIQUE,
  -- The user's password hash.
  passhash bytea NOT NULL,
  -- The user's display name.
  name text NOT NULL,
  -- The user's locale.
  locale locale NOT NULL DEFAULT 'en-US',
  -- The time the user was created.
  registered_at timestamp NOT NULL DEFAULT now(),
  -- The notification service for this user. If null, notifications are
  -- disabled.
  notification_service notificationconfig,
  -- The custom notification message for this user. If null, the default
  -- message is used.
  custom_notification notification
);

CREATE TABLE user_sessions (
  -- The session ID. This should never be used to log in, but it can be used
  -- to revoke a session.
  session_id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  -- The user's ID in xid format.
  user_id xid NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
  -- The session token.
  token bytea UNIQUE NOT NULL,
  -- The time the session was created.
  created_at timestamp NOT NULL DEFAULT now(),
  -- The time the session was last used.
  last_used timestamp NOT NULL DEFAULT now(),
  -- The user agent string, if any. Collected for management purposes.
  user_agent text
);

CREATE TABLE user_avatars (
  -- The user's ID in xid format.
  user_id xid PRIMARY KEY REFERENCES users (user_id) ON DELETE CASCADE,
  -- The MIME type of the image.
  mime_type text NOT NULL,
  -- The compression of the image data, if any.
  compression compression,
  -- The user's avatar image, limited to 1MB.
  avatar_image bytea NOT NULL CHECK (octet_length(avatar_image) <= 1048576)
);

CREATE TABLE dosage_schedule (
  user_id xid PRIMARY KEY REFERENCES users (user_id) ON DELETE CASCADE,
  -- The delivery method of the medication.
  delivery_method text REFERENCES delivery_methods (id) ON DELETE CASCADE,
  -- The dose of the medication.
  -- The units are determined by the delivery method.
  dose numeric(4, 2) NOT NULL,
  -- The interval between doses.
  interval interval NOT NULL CHECK (interval > '0 minutes'::interval),
  -- How many patches are on at a time. Only relevant (non-null) for patches.
  concurrence numeric(2, 0)
);

CREATE TABLE dosage_history (
  dose_id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  -- The dose that is previous to this one.
  -- This is mostly used for reconciling dose conflicts.
  last_dose bigint REFERENCES dosage_history (dose_id) ON DELETE SET NULL,
  -- The user that took the dose.
  user_id xid REFERENCES users (user_id) ON DELETE CASCADE,
  -- The delivery method of the medication.
  delivery_method text REFERENCES delivery_methods (id) ON DELETE CASCADE,
  -- The dose of the medication.
  -- This is usually copied from the schedule, but can be overridden.
  dose numeric(4, 2) NOT NULL,
  -- The time the dose was taken.
  taken_at timestamptz NOT NULL,
  -- The time the patch was taken off. This is only applicable to patches.
  taken_off_at timestamptz
);

-- Ensure a binary-tree index on user_id and taken_at so that we can query based
-- on the timestamp.
CREATE INDEX dosage_history_user_id_taken_at ON dosage_history USING btree (user_id, taken_at);

CREATE TABLE notification_history (
  notification_id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  -- The user that the notification was for.
  user_id xid REFERENCES users (user_id) ON DELETE CASCADE,
  -- The dosage that the notification was for.
  dosage_id bigint REFERENCES dosage_history (dose_id) ON DELETE CASCADE,
  -- The time the notification was sent.
  sent_at timestamptz NOT NULL,
  -- The error if the notification failed to send.
  error_reason text
);
