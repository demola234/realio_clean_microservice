// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type Sessions struct {
	// Unique identifier for each session.
	SessionID uuid.UUID `json:"session_id"`
	// Foreign key linking to the user table (identifies the user).
	UserID uuid.UUID `json:"user_id"`
	// The session token, which can be a JWT or another token format.
	Token        string         `json:"token"`
	Otp          sql.NullString `json:"otp"`
	OtpExpiresAt sql.NullTime   `json:"otp_expires_at"`
	OtpAttempts  sql.NullInt32  `json:"otp_attempts"`
	OtpVerified  sql.NullBool   `json:"otp_verified"`
	// Timestamp of when the session was created.
	CreatedAt time.Time `json:"created_at"`
	// Timestamp of when the session expires.
	ExpiresAt time.Time `json:"expires_at"`
	// Tracks the last activity time for session timeout checks.
	LastActivity time.Time `json:"last_activity"`
	// The IP address from which the session was initiated.
	IpAddress sql.NullString `json:"ip_address"`
	// The user agent (browser or device info) for the session.
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
	Name           string         `json:"name"`
	Username       string         `json:"username"`
	ProfilePicture sql.NullString `json:"profile_picture"`
	Bio            sql.NullString `json:"bio"`
	// User's email (unique)
	Email string `json:"email"`
	// Hashed password
	Password string `json:"password"`
	// Role
	Role sql.NullString `json:"role"`
	// Contact number
	Phone sql.NullString `json:"phone"`
	// Timestamp of registration
	CreatedAt sql.NullTime `json:"created_at"`
	// Timestamp of last update
	UpdatedAt sql.NullTime `json:"updated_at"`
}
