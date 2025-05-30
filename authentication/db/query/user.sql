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
    email_verified,
    is_active,
    last_login,
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
    $12, -- email_verified
    $13, -- is_active
    $14, -- last_login
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

-- name: UpdateUserProfilePicture :one
UPDATE users
SET profile_picture = $2,
    updated_at = now()
    WHERE id = $1
    RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateLastLogin :exec
UPDATE users
SET last_login = now()
WHERE id = $1;

-- name: UpdateEmailVerification :exec
UPDATE users
SET email_verified = true
WHERE id = $1;

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