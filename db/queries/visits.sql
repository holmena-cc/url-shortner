-- name: CreateVisit :one
INSERT INTO visits (url_id, ip_address, referrer, country)
VALUES ($1, $2, $3, $4)
RETURNING click_id, url_id, click_date, ip_address, referrer, country;

-- name: GetVisitByID :one
SELECT click_id, url_id, click_date, ip_address, referrer, country
FROM visits
WHERE click_id = $1;

-- name: ListVisitsByURL :many
SELECT click_id, click_date, ip_address, referrer, country
FROM visits
WHERE url_id = $1
ORDER BY click_date DESC;

-- name: DeleteVisit :exec
DELETE FROM visits WHERE click_id = $1;
