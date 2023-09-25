-- name: CreateUser :one
INSERT INTO users (
  name, bio, birth_date, last_location_long, last_location_lat,
  mobile, last_active, email, sex, interested_in
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;