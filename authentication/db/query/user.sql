-- name: CreateUser :one
INSERT INTO
    users (
        email,
        hashed_password,
        full_name,
        role
    )
VALUES (
        $1, -- email
        $2, -- hash_password
        $3, -- role
        $4 -- full_name
    ) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    full_name = COALESCE($1, full_name),
    hashed_password = COALESCE($2, hashed_password),
    email = COALESCE($3, email),
    role = COALESCE($4, role),
    updated_at = now()
WHERE
    email = $5 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE email = $1;