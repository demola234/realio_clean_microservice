
-- name: CreatePasswordReset :one
INSERT INTO password_resets (
    id, user_id, token, expires_at, created_at, used
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPasswordResetByToken :one
SELECT * FROM password_resets
WHERE token = $1 AND used = false AND expires_at > now()
LIMIT 1;

-- name: InvalidatePasswordReset :one
UPDATE password_resets
SET used = true, updated_at = now()
WHERE token = $1 AND used = false
RETURNING *;

-- name: DeletePasswordResetsByUserId :exec
DELETE FROM password_resets
WHERE user_id = $1;

-- name: DeleteExpiredPasswordResets :exec
DELETE FROM password_resets
WHERE expires_at < now() - INTERVAL '7 days';