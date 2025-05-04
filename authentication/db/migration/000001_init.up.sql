CREATE TABLE "users" (
    "id" UUID UNIQUE PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "username" VARCHAR UNIQUE NOT NULL,
    "profile_picture" VARCHAR,
    "bio" VARCHAR,
    "email" VARCHAR UNIQUE NOT NULL,
    "password" VARCHAR,
    "role" VARCHAR,
    "phone" VARCHAR,
    "provider" VARCHAR DEFAULT 'local',
    "provider_id" VARCHAR UNIQUE,
    "email_verified" BOOLEAN DEFAULT false,
    "is_active" BOOLEAN DEFAULT true,
    "last_login" TIMESTAMP,
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
COMMENT ON COLUMN "users"."username" IS 'Unique username';
COMMENT ON COLUMN "users"."profile_picture" IS 'URL to profile picture';
COMMENT ON COLUMN "users"."bio" IS 'Short user bio';
COMMENT ON COLUMN "users"."email" IS 'User''s email (unique)';
COMMENT ON COLUMN "users"."password" IS 'Hashed password (nullable for OAuth users)';
COMMENT ON COLUMN "users"."role" IS 'Role (admin, user, etc.)';
COMMENT ON COLUMN "users"."phone" IS 'Contact number';
COMMENT ON COLUMN "users"."provider" IS 'Authentication provider (local, google, github, etc.)';
COMMENT ON COLUMN "users"."provider_id" IS 'User ID from OAuth provider (unique)';
COMMENT ON COLUMN "users"."created_at" IS 'Timestamp of registration';
COMMENT ON COLUMN "users"."updated_at" IS 'Timestamp of last update';
COMMENT ON COLUMN "users"."email_verified" IS 'Indicates if email is verified';
COMMENT ON COLUMN "users"."is_active" IS 'Indicates if user is active';
COMMENT ON COLUMN "users"."last_login" IS 'Timestamp of last login';

-- Comments for sessions table
COMMENT ON COLUMN "sessions"."session_id" IS 'Unique identifier for each session.';
COMMENT ON COLUMN "sessions"."user_id" IS 'Foreign key linking to the user table (identifies the user).';
COMMENT ON COLUMN "sessions"."token" IS 'The session token (JWT or other).';
COMMENT ON COLUMN "sessions"."otp" IS 'One-time password (for MFA)';
COMMENT ON COLUMN "sessions"."otp_expires_at" IS 'Expiry time for OTP';
COMMENT ON COLUMN "sessions"."otp_attempts" IS 'Number of OTP attempts made';
COMMENT ON COLUMN "sessions"."otp_verified" IS 'Indicates if OTP was verified';
COMMENT ON COLUMN "sessions"."created_at" IS 'Timestamp of when the session was created.';
COMMENT ON COLUMN "sessions"."expires_at" IS 'Timestamp of when the session expires.';
COMMENT ON COLUMN "sessions"."last_activity" IS 'Tracks the last activity time for session timeout checks.';
COMMENT ON COLUMN "sessions"."ip_address" IS 'The IP address from which the session was initiated.';
COMMENT ON COLUMN "sessions"."user_agent" IS 'The user agent (browser/device info) for the session.';
COMMENT ON COLUMN "sessions"."is_active" IS 'Indicates whether the session is currently active.';
COMMENT ON COLUMN "sessions"."revoked_at" IS 'Timestamp for when the session was revoked, if applicable.';
COMMENT ON COLUMN "sessions"."device_info" IS 'Stores additional device details if needed.';

-- Add foreign key constraint
ALTER TABLE "sessions"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

-- Add password/provider constraint
ALTER TABLE "users"
ADD CONSTRAINT users_password_or_provider_id_check
CHECK (
    (provider = 'local' AND password IS NOT NULL)
    OR
    (provider != 'local' AND provider_id IS NOT NULL)
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON "users"("email");
CREATE INDEX idx_users_username ON "users"("username");
CREATE INDEX idx_users_provider_id ON "users"("provider_id") WHERE "provider_id" IS NOT NULL;
CREATE INDEX idx_users_is_active ON "users"("is_active") WHERE "is_active" = true;
CREATE INDEX idx_sessions_user_id ON "sessions"("user_id");
CREATE INDEX idx_sessions_token ON "sessions"("token");
CREATE INDEX idx_sessions_is_active ON "sessions"("is_active") WHERE "is_active" = true;
CREATE INDEX idx_sessions_expires_at ON "sessions"("expires_at");

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for updating updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON "users"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();