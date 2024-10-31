-- name: CreateSession :one
INSERT INTO
    sessions (
        session_id,
        user_id,
        token,
        created_at,
        expires_at,
        last_activity,
        ip_address,
        user_agent,
        is_active,
        revoked_at,
        otp_verified,
        otp_expires_at,
        otp_attempts,
        otp,
        device_info
    )
VALUES (
        $1, -- session_id (unique identifier for the session, UUID)
        $2, -- user_id (UUID of the associated user)
        $3, -- token (session or refresh token for authentication)
        now(), -- created_at (timestamp for when session is created)
        $4, -- expires_at (timestamp for when session expires)
        $5, -- last_activity (timestamp for last user activity)
        $6, -- ip_address (IP address of the client)
        $7, -- user_agent (device or browser information)
        $8, -- is_active (boolean to indicate if the session is active)
        $9, -- revoked_at (timestamp for when session is revoked, if applicable)
        $10,
        $11,
        $12,
        $13,
        $14
    ) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions 
WHERE session_id = $1 
   OR user_id = $1;

-- name: UpdateSessionActivity :exec
UPDATE sessions
SET
    last_activity = $1,
    is_active = $2
WHERE
    session_id = $3;

-- name: RevokeSession :exec
UPDATE sessions
SET
    is_active = false,
    revoked_at = now()
WHERE
    session_id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE session_id = $1;

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
WHERE
    user_id = $13 RETURNING *;