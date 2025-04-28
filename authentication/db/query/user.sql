-- name: CreateUser :one
INSERT INTO users (
    id,
    name,
    username,
    email,
    password,
    profile_picture,
    bio,
    role,
    phone,
    provider,
    provider_id,
    created_at,
    updated_at
) VALUES (
    $1, -- id
    $2, -- name
    $3, -- username
    $4, -- email
    $5, -- password
    $6, -- profile_picture
    $7, -- bio
    $8, -- role
    $9, -- phone
    $10, -- provider
    $11, -- provider_id
    now(), -- created_at
    now()  -- updated_at
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 OR id::text = $1 OR username = $1
LIMIT 1;

-- name: CheckEmailExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE email = $1
    LIMIT 1
) AS exists;

-- name: ChangePassword :one
UPDATE users
SET password = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    name = COALESCE($1, name),
    username = COALESCE($2, username),
    email = COALESCE($3, email),
    password = COALESCE($4, password),
    profile_picture = COALESCE($5, profile_picture),
    bio = COALESCE($6, bio),
    role = COALESCE($7, role),
    phone = COALESCE($8, phone),
    updated_at = now()
WHERE id = $9
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;