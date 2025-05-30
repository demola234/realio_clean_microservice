// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type PasswordResets struct {
	ID        uuid.UUID    `json:"id"`
	UserID    uuid.UUID    `json:"user_id"`
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	Used      bool         `json:"used"`
}

type Sessions struct {
	// Unique identifier for each session.
	SessionID uuid.UUID `json:"session_id"`
	// Foreign key linking to the user table (identifies the user).
	UserID uuid.UUID `json:"user_id"`
	// The session token (JWT or other).
	Token string `json:"token"`
	// One-time password (for MFA)
	Otp sql.NullString `json:"otp"`
	// Expiry time for OTP
	OtpExpiresAt sql.NullTime `json:"otp_expires_at"`
	// Number of OTP attempts made
	OtpAttempts sql.NullInt32 `json:"otp_attempts"`
	// Indicates if OTP was verified
	OtpVerified sql.NullBool `json:"otp_verified"`
	// Timestamp of when the session was created.
	CreatedAt time.Time `json:"created_at"`
	// Timestamp of when the session expires.
	ExpiresAt time.Time `json:"expires_at"`
	// Tracks the last activity time for session timeout checks.
	LastActivity time.Time `json:"last_activity"`
	// The IP address from which the session was initiated.
	IpAddress sql.NullString `json:"ip_address"`
	// The user agent (browser/device info) for the session.
	UserAgent sql.NullString `json:"user_agent"`
	// Indicates whether the session is currently active.
	IsActive bool `json:"is_active"`
	// Timestamp for when the session was revoked, if applicable.
	RevokedAt sql.NullTime `json:"revoked_at"`
	// Stores additional device details if needed.
	DeviceInfo pqtype.NullRawMessage `json:"device_info"`
}

type Users struct {
	// Primary key
	ID uuid.UUID `json:"id"`
	// User's full name
	Name string `json:"name"`
	// Unique username
	Username string `json:"username"`
	// URL to profile picture
	ProfilePicture sql.NullString `json:"profile_picture"`
	// Short user bio
	Bio sql.NullString `json:"bio"`
	// User's email (unique)
	Email string `json:"email"`
	// Hashed password (nullable for OAuth users)
	Password sql.NullString `json:"password"`
	// Role (admin, user, etc.)
	Role sql.NullString `json:"role"`
	// Contact number
	Phone sql.NullString `json:"phone"`
	// Authentication provider (local, google, github, etc.)
	Provider sql.NullString `json:"provider"`
	// User ID from OAuth provider (unique)
	ProviderID sql.NullString `json:"provider_id"`
	// Indicates if email is verified
	EmailVerified sql.NullBool `json:"email_verified"`
	// Indicates if user is active
	IsActive sql.NullBool `json:"is_active"`
	// Timestamp of last login
	LastLogin sql.NullTime `json:"last_login"`
	// Timestamp of registration
	CreatedAt sql.NullTime `json:"created_at"`
	// Timestamp of last update
	UpdatedAt sql.NullTime `json:"updated_at"`
}
