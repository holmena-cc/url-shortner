-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING user_id, username, email, register_date;

-- name: GetUserByID :one
SELECT user_id, username, email, password_hash, register_date
FROM users
WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT user_id, username, email, password_hash, register_date
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT user_id, username, email, register_date
FROM users
ORDER BY register_date DESC;

-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3
WHERE user_id = $1
RETURNING user_id, username, email, register_date;

-- name: DeleteUser :exec
DELETE FROM users WHERE user_id = $1;
