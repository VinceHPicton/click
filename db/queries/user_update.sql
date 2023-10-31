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