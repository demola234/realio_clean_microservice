-- name: CreateSession :one
INSERT INTO sessions (
    session_id,
    user_id,
    token,
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
    $4, -- expires_at
    $5, -- last_activity
    $6, -- ip_address
    $7, -- user_agent
    $8, -- is_active
    $9, -- revoked_at
    $10 -- device_info
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
    expires_at = COALESCE($2, expires_at),
    last_activity = COALESCE($3, last_activity),
    ip_address = COALESCE($4, ip_address),
    user_agent = COALESCE($5, user_agent),
    is_active = COALESCE($6, is_active),
    revoked_at = COALESCE($7, revoked_at),
    device_info = COALESCE($8, device_info)
WHERE user_id = $9
RETURNING *;

-- name: CreateLoginHistoryEntry :one
INSERT INTO sessions (
    session_id, user_id, ip_address, user_agent
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetSessionsByUserID :many
SELECT * FROM sessions 
WHERE user_id = $1 
ORDER BY created_at DESC;
