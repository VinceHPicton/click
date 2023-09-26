-- name: CreateUser :one
INSERT INTO users (
  name, bio, birth_date, mobile, email, sex, interested_in
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;