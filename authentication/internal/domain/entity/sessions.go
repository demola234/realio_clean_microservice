package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// Session entity based on the sessions table schema
type Session struct {
	SessionID    uuid.UUID              `json:"session_id"`
	UserID       uuid.UUID              `json:"user_id"`
	Token        string                 `json:"token"`
	Otp          string                 `json:"otp"`
	OtpExpiresAt time.Time              `json:"otp_expires_at"`
	OtpAttempts  int                    `json:"otp_attempts"`
	OTPVerified  bool                   `json:"otp_verified"`
	CreatedAt    time.Time              `json:"created_at"`
	ExpiresAt    time.Time              `json:"expires_at"`
	LastActivity time.Time              `json:"last_activity"`
	IpAddress    string                 `json:"ip_address,omitempty"`
	UserAgent    string                 `json:"user_agent,omitempty"`
	IsActive     bool                   `json:"is_active"`
	RevokedAt    *time.Time             `json:"revoked_at,omitempty"`
	DeviceInfo   *pqtype.NullRawMessage `json:"device_info,omitempty"`
}

type UpdateOtp struct {
	Otp          string    `json:"otp"`
	OtpExpiresAt time.Time `json:"otp_expires_at"`
	OtpAttempts  int       `json:"otp_attempts"`
	Email        string    `json:"email"`
	OTPVerified  bool      `json:"otp_verified"`
}
