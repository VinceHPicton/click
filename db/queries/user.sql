-- name: CreateUser :one

INSERT INTO app.users (
  name, bio, birth_date, mobile, email, sex, interested_in
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: HardDeleteUser :exec
DELETE FROM app.users WHERE id = $1;

-- name: SoftDeleteUser :exec
UPDATE app.users
  set deleted_at = NOW()
WHERE id = $1;

-- name: GetUser :one
SELECT * FROM app.users
WHERE id = $1 LIMIT 1;

-- name: GetActiveUserByMobile :many
SELECT * FROM app.users
WHERE mobile = $1 AND deleted_at IS NULL AND banned_at IS NULL;

-- name: BanUser :exec
UPDATE app.users
  set banned_at = NOW()
WHERE id = $1;

-- name: UpdateUser :one
UPDATE app.users
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

-- name: GetAllUsers :many
SELECT * FROM app.users;