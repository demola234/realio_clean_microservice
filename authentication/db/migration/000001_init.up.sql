CREATE TABLE "users" (
    "id" UUID UNIQUE PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "username" VARCHAR UNIQUE NOT NULL,
    "profile_picture" VARCHAR,
    "bio" VARCHAR,
    "email" VARCHAR UNIQUE NOT NULL,
    "password" VARCHAR NOT NULL,
    "role" VARCHAR,
    "phone" VARCHAR,
    "created_at" TIMESTAMP DEFAULT now(),
    "updated_at" TIMESTAMP DEFAULT now()
);

CREATE TABLE "sessions" (
    "session_id" UUID UNIQUE PRIMARY KEY,
    "user_id" UUID NOT NULL,
    "token" VARCHAR(255) NOT NULL,
    "otp" VARCHAR(6),
    "otp_expires_at" TIMESTAMP,
    "otp_attempts" INT DEFAULT 0,
    "otp_verified" BOOLEAN DEFAULT false,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "expires_at" TIMESTAMP NOT NULL,
    "last_activity" TIMESTAMP NOT NULL,
    "ip_address" VARCHAR(45),
    "user_agent" VARCHAR(255),
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "revoked_at" TIMESTAMP,
    "device_info" JSON
);

COMMENT ON COLUMN "users"."id" IS 'Primary key';

COMMENT ON COLUMN "users"."name" IS 'User''s full name';

COMMENT ON COLUMN "users"."email" IS 'User''s email (unique)';

COMMENT ON COLUMN "users"."password" IS 'Hashed password';

COMMENT ON COLUMN "users"."role" IS 'Role';

COMMENT ON COLUMN "users"."phone" IS 'Contact number';

COMMENT ON COLUMN "users"."created_at" IS 'Timestamp of registration';

COMMENT ON COLUMN "users"."updated_at" IS 'Timestamp of last update';

COMMENT ON COLUMN "sessions"."session_id" IS 'Unique identifier for each session.';

COMMENT ON COLUMN "sessions"."user_id" IS 'Foreign key linking to the user table (identifies the user).';

COMMENT ON COLUMN "sessions"."token" IS 'The session token, which can be a JWT or another token format.';

COMMENT ON COLUMN "sessions"."created_at" IS 'Timestamp of when the session was created.';

COMMENT ON COLUMN "sessions"."expires_at" IS 'Timestamp of when the session expires.';

COMMENT ON COLUMN "sessions"."last_activity" IS 'Tracks the last activity time for session timeout checks.';

COMMENT ON COLUMN "sessions"."ip_address" IS 'The IP address from which the session was initiated.';

COMMENT ON COLUMN "sessions"."user_agent" IS 'The user agent (browser or device info) for the session.';

COMMENT ON COLUMN "sessions"."is_active" IS 'Indicates whether the session is currently active.';

COMMENT ON COLUMN "sessions"."revoked_at" IS 'Timestamp for when the session was revoked, if applicable.';

COMMENT ON COLUMN "sessions"."device_info" IS 'Stores additional device details if needed.';

ALTER TABLE "sessions"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");