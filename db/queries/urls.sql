-- name: CreateURL :one
INSERT INTO urls (user_id, original_url, short_code, custom_alias)
VALUES ($1, $2, $3, $4)
RETURNING url_id, user_id, original_url, short_code, custom_alias, creation_date;

-- name: GetURLByID :one
SELECT url_id, user_id, original_url, short_code, custom_alias, creation_date
FROM urls
WHERE url_id = $1;

-- name: GetURLByCode :one
SELECT url_id, user_id, original_url, short_code, custom_alias, creation_date
FROM urls
WHERE short_code = $1;

-- name: ListURLsByUser :many
SELECT url_id, short_code, original_url, custom_alias, creation_date
FROM urls
WHERE user_id = $1
ORDER BY creation_date DESC;

-- name: UpdateURL :one
UPDATE urls
SET original_url = $2, short_code = $3, custom_alias = $4
WHERE url_id = $1
RETURNING url_id, short_code, original_url, custom_alias, creation_date;

-- name: DeleteURL :exec
DELETE FROM urls WHERE url_id = $1;

-- name: AliasExists :one
SELECT EXISTS (
    SELECT 1
    FROM urls
    WHERE custom_alias = $1
);