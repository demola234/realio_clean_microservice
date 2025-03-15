-- name: CreateUser :one
INSERT INTO
    users (
        id,
        name,
        email,
        password,
        profile_picture,
        username,
        bio,
        role,
        phone,
        created_at,
        updated_at
    )
VALUES (
        $1, -- id (UUID for the user)
        $2, -- name (user's full name)
        $3, -- email (user's unique email)
        $4, -- password (hashed password)
        $5, -- role (user role, default is 'buyer')
        $6, -- phone (optional phone number)
        $7, -- bio (optional user bio)
        $8, --profile_picture
        $9, -- username (unique username)
        now(), -- created_at (timestamp for user creation)
        now() -- updated_at (timestamp for the last update)
    ) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1 OR id::text = $1 LIMIT 1;

-- name: ChangePassword :one
UPDATE users
SET password = $2,
updated_at = now()
WHERE
    id = $1 RETURNING *;

-- name: CheckEmailExists :one
SELECT EXISTS (
        SELECT 1
        FROM users
        WHERE
            email = $1
        LIMIT 1
    ) AS exists;

-- name: UpdateUser :one
UPDATE users
SET
    name = COALESCE($1, name),
    email = COALESCE($2, email),
    password = COALESCE($3, password),
    role = COALESCE($4, role),
    phone = COALESCE($5, phone),
    updated_at = now()
WHERE
    id = $6 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;