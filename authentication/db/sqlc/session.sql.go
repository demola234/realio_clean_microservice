// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: session.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

const createLoginHistoryEntry = `-- name: CreateLoginHistoryEntry :one
INSERT INTO sessions (
    session_id, user_id, ip_address, user_agent
) VALUES (
    $1, $2, $3, $4
) RETURNING session_id, user_id, token, otp, otp_expires_at, otp_attempts, otp_verified, created_at, expires_at, last_activity, ip_address, user_agent, is_active, revoked_at, device_info
`

type CreateLoginHistoryEntryParams struct {
	SessionID uuid.UUID      `json:"session_id"`
	UserID    uuid.UUID      `json:"user_id"`
	IpAddress sql.NullString `json:"ip_address"`
	UserAgent sql.NullString `json:"user_agent"`
}

func (q *Queries) CreateLoginHistoryEntry(ctx context.Context, arg CreateLoginHistoryEntryParams) (Sessions, error) {
	row := q.db.QueryRowContext(ctx, createLoginHistoryEntry,
		arg.SessionID,
		arg.UserID,
		arg.IpAddress,
		arg.UserAgent,
	)
	var i Sessions
	err := row.Scan(
		&i.SessionID,
		&i.UserID,
		&i.Token,
		&i.Otp,
		&i.OtpExpiresAt,
		&i.OtpAttempts,
		&i.OtpVerified,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.LastActivity,
		&i.IpAddress,
		&i.UserAgent,
		&i.IsActive,
		&i.RevokedAt,
		&i.DeviceInfo,
	)
	return i, err
}

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (
    session_id,
    user_id,
    token,
    otp,
    otp_expires_at,
    otp_attempts,
    otp_verified,
    created_at,
    expires_at,
    last_activity,
    ip_address,
    user_agent,
    is_active,
    revoked_at,
    device_info
) VALUES (
    $1, -- session_id
    $2, -- user_id
    $3, -- token
    $4, -- otp
    $5, -- otp_expires_at
    $6, -- otp_attempts
    $7, -- otp_verified
    now(), -- created_at
    $8, -- expires_at
    $9, -- last_activity
    $10, -- ip_address
    $11, -- user_agent
    $12, -- is_active
    $13, -- revoked_at
    $14 -- device_info
) RETURNING session_id, user_id, token, otp, otp_expires_at, otp_attempts, otp_verified, created_at, expires_at, last_activity, ip_address, user_agent, is_active, revoked_at, device_info
`

type CreateSessionParams struct {
	SessionID    uuid.UUID             `json:"session_id"`
	UserID       uuid.UUID             `json:"user_id"`
	Token        string                `json:"token"`
	Otp          sql.NullString        `json:"otp"`
	OtpExpiresAt sql.NullTime          `json:"otp_expires_at"`
	OtpAttempts  sql.NullInt32         `json:"otp_attempts"`
	OtpVerified  sql.NullBool          `json:"otp_verified"`
	ExpiresAt    time.Time             `json:"expires_at"`
	LastActivity time.Time             `json:"last_activity"`
	IpAddress    sql.NullString        `json:"ip_address"`
	UserAgent    sql.NullString        `json:"user_agent"`
	IsActive     bool                  `json:"is_active"`
	RevokedAt    sql.NullTime          `json:"revoked_at"`
	DeviceInfo   pqtype.NullRawMessage `json:"device_info"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Sessions, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.SessionID,
		arg.UserID,
		arg.Token,
		arg.Otp,
		arg.OtpExpiresAt,
		arg.OtpAttempts,
		arg.OtpVerified,
		arg.ExpiresAt,
		arg.LastActivity,
		arg.IpAddress,
		arg.UserAgent,
		arg.IsActive,
		arg.RevokedAt,
		arg.DeviceInfo,
	)
	var i Sessions
	err := row.Scan(
		&i.SessionID,
		&i.UserID,
		&i.Token,
		&i.Otp,
		&i.OtpExpiresAt,
		&i.OtpAttempts,
		&i.OtpVerified,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.LastActivity,
		&i.IpAddress,
		&i.UserAgent,
		&i.IsActive,
		&i.RevokedAt,
		&i.DeviceInfo,
	)
	return i, err
}

const deleteSession = `-- name: DeleteSession :exec
DELETE FROM sessions
WHERE session_id = $1
`

func (q *Queries) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteSession, sessionID)
	return err
}

const getLoginHistory = `-- name: GetLoginHistory :many
SELECT session_id, user_id, token, otp, otp_expires_at, otp_attempts, otp_verified, created_at, expires_at, last_activity, ip_address, user_agent, is_active, revoked_at, device_info FROM sessions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
`

type GetLoginHistoryParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
}

func (q *Queries) GetLoginHistory(ctx context.Context, arg GetLoginHistoryParams) ([]Sessions, error) {
	rows, err := q.db.QueryContext(ctx, getLoginHistory, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Sessions{}
	for rows.Next() {
		var i Sessions
		if err := rows.Scan(
			&i.SessionID,
			&i.UserID,
			&i.Token,
			&i.Otp,
			&i.OtpExpiresAt,
			&i.OtpAttempts,
			&i.OtpVerified,
			&i.CreatedAt,
			&i.ExpiresAt,
			&i.LastActivity,
			&i.IpAddress,
			&i.UserAgent,
			&i.IsActive,
			&i.RevokedAt,
			&i.DeviceInfo,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSessionByID = `-- name: GetSessionByID :one
SELECT session_id, user_id, token, otp, otp_expires_at, otp_attempts, otp_verified, created_at, expires_at, last_activity, ip_address, user_agent, is_active, revoked_at, device_info FROM sessions
WHERE session_id = $1
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (Sessions, error) {
	row := q.db.QueryRowContext(ctx, getSessionByID, sessionID)
	var i Sessions
	err := row.Scan(
		&i.SessionID,
		&i.UserID,
		&i.Token,
		&i.Otp,
		&i.OtpExpiresAt,
		&i.OtpAttempts,
		&i.OtpVerified,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.LastActivity,
		&i.IpAddress,
		&i.UserAgent,
		&i.IsActive,
		&i.RevokedAt,
		&i.DeviceInfo,
	)
	return i, err
}

const getSessionByUserID = `-- name: GetSessionByUserID :one
SELECT session_id, user_id, token, otp, otp_expires_at, otp_attempts, otp_verified, created_at, expires_at, last_activity, ip_address, user_agent, is_active, revoked_at, device_info FROM sessions
WHERE user_id = $1
LIMIT 1
`

func (q *Queries) GetSessionByUserID(ctx context.Context, userID uuid.UUID) (Sessions, error) {
	row := q.db.QueryRowContext(ctx, getSessionByUserID, userID)
	var i Sessions
	err := row.Scan(
		&i.SessionID,
		&i.UserID,
		&i.Token,
		&i.Otp,
		&i.OtpExpiresAt,
		&i.OtpAttempts,
		&i.OtpVerified,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.LastActivity,
		&i.IpAddress,
		&i.UserAgent,
		&i.IsActive,
		&i.RevokedAt,
		&i.DeviceInfo,
	)
	return i, err
}

const getSessionsByUserID = `-- name: GetSessionsByUserID :many
SELECT session_id, user_id, token, otp, otp_expires_at, otp_attempts, otp_verified, created_at, expires_at, last_activity, ip_address, user_agent, is_active, revoked_at, device_info FROM sessions 
WHERE user_id = $1 
ORDER BY created_at DESC
`

func (q *Queries) GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]Sessions, error) {
	rows, err := q.db.QueryContext(ctx, getSessionsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Sessions{}
	for rows.Next() {
		var i Sessions
		if err := rows.Scan(
			&i.SessionID,
			&i.UserID,
			&i.Token,
			&i.Otp,
			&i.OtpExpiresAt,
			&i.OtpAttempts,
			&i.OtpVerified,
			&i.CreatedAt,
			&i.ExpiresAt,
			&i.LastActivity,
			&i.IpAddress,
			&i.UserAgent,
			&i.IsActive,
			&i.RevokedAt,
			&i.DeviceInfo,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const revokeSession = `-- name: RevokeSession :exec
UPDATE sessions
SET
    is_active = false,
    revoked_at = now()
WHERE user_id = $1 AND is_active = true
`

func (q *Queries) RevokeSession(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, revokeSession, userID)
	return err
}

const updateSession = `-- name: UpdateSession :one
UPDATE sessions
SET
    token = COALESCE($1, token),
    otp = COALESCE($2, otp),
    otp_expires_at = COALESCE($3, otp_expires_at),
    otp_attempts = COALESCE($4, otp_attempts),
    expires_at = COALESCE($5, expires_at),
    last_activity = COALESCE($6, last_activity),
    ip_address = COALESCE($7, ip_address),
    user_agent = COALESCE($8, user_agent),
    is_active = COALESCE($9, is_active),
    revoked_at = COALESCE($10, revoked_at),
    otp_verified = COALESCE($11, otp_verified),
    device_info = COALESCE($12, device_info)
WHERE user_id = $13
RETURNING session_id, user_id, token, otp, otp_expires_at, otp_attempts, otp_verified, created_at, expires_at, last_activity, ip_address, user_agent, is_active, revoked_at, device_info
`

type UpdateSessionParams struct {
	Token        string                `json:"token"`
	Otp          sql.NullString        `json:"otp"`
	OtpExpiresAt sql.NullTime          `json:"otp_expires_at"`
	OtpAttempts  sql.NullInt32         `json:"otp_attempts"`
	ExpiresAt    time.Time             `json:"expires_at"`
	LastActivity time.Time             `json:"last_activity"`
	IpAddress    sql.NullString        `json:"ip_address"`
	UserAgent    sql.NullString        `json:"user_agent"`
	IsActive     bool                  `json:"is_active"`
	RevokedAt    sql.NullTime          `json:"revoked_at"`
	OtpVerified  sql.NullBool          `json:"otp_verified"`
	DeviceInfo   pqtype.NullRawMessage `json:"device_info"`
	UserID       uuid.UUID             `json:"user_id"`
}

func (q *Queries) UpdateSession(ctx context.Context, arg UpdateSessionParams) (Sessions, error) {
	row := q.db.QueryRowContext(ctx, updateSession,
		arg.Token,
		arg.Otp,
		arg.OtpExpiresAt,
		arg.OtpAttempts,
		arg.ExpiresAt,
		arg.LastActivity,
		arg.IpAddress,
		arg.UserAgent,
		arg.IsActive,
		arg.RevokedAt,
		arg.OtpVerified,
		arg.DeviceInfo,
		arg.UserID,
	)
	var i Sessions
	err := row.Scan(
		&i.SessionID,
		&i.UserID,
		&i.Token,
		&i.Otp,
		&i.OtpExpiresAt,
		&i.OtpAttempts,
		&i.OtpVerified,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.LastActivity,
		&i.IpAddress,
		&i.UserAgent,
		&i.IsActive,
		&i.RevokedAt,
		&i.DeviceInfo,
	)
	return i, err
}

const updateSessionActivity = `-- name: UpdateSessionActivity :exec
UPDATE sessions
SET
    last_activity = $1,
    is_active = $2
WHERE session_id = $3
`

type UpdateSessionActivityParams struct {
	LastActivity time.Time `json:"last_activity"`
	IsActive     bool      `json:"is_active"`
	SessionID    uuid.UUID `json:"session_id"`
}

func (q *Queries) UpdateSessionActivity(ctx context.Context, arg UpdateSessionActivityParams) error {
	_, err := q.db.ExecContext(ctx, updateSessionActivity, arg.LastActivity, arg.IsActive, arg.SessionID)
	return err
}
