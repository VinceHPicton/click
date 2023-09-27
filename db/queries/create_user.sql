-- name: CreateUser :one
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO users (
  name, bio, birth_date, mobile, email, sex, interested_in
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;