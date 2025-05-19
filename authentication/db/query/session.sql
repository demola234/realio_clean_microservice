-- name: CreateSession :one
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
) RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE session_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: GetSessionByUserID :one
SELECT * FROM sessions
WHERE user_id = $1
LIMIT 1;

-- name: UpdateSessionActivity :exec
UPDATE sessions
SET
    last_activity = $1,
    is_active = $2
WHERE session_id = $3;

-- name: RevokeSession :exec
UPDATE sessions
SET
    is_active = false,
    revoked_at = now()
WHERE user_id = $1 AND is_active = true;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE session_id = $1;

-- name: UpdateSession :one
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
RETURNING *;

-- name: CreateLoginHistoryEntry :one
INSERT INTO sessions (
    session_id, user_id, ip_address, user_agent
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetLoginHistory :many
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetSessionsByUserID :many
SELECT * FROM sessions 
WHERE user_id = $1 
ORDER BY created_at DESC;
