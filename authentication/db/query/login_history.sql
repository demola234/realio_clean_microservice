-- name: CreateLoginHistory :one
INSERT INTO login_history (
    id, user_id, ip_address, user_agent
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetLoginHistory :many
SELECT * FROM login_history
WHERE user_id = $1
ORDER BY login_time DESC
LIMIT $2;
