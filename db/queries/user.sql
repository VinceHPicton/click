-- name: UserCreate :one

INSERT INTO users (
  name, bio, birth_date, mobile, email, sex, interested_in
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UserDelete :exec
DELETE FROM users WHERE id = $1;

-- name: UserGet :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UserUpdate :one
UPDATE users
  set
  name = $8,
  bio = $7,
  birth_date = $6,
  mobile = $5,
  email = $4,
  sex = $3,
  interested_in = $2
WHERE id = $1
RETURNING *;