-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING user_id, email, register_date;

-- name: GetUserByID :one
SELECT user_id, email, password_hash, register_date
FROM users
WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT user_id, email, password_hash, register_date
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT user_id, email, register_date
FROM users
ORDER BY register_date DESC;

-- name: DeleteUser :exec
DELETE FROM users WHERE user_id = $1;
