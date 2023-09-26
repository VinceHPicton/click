-- name: UpdateUser :one
UPDATE users
  set name = $11,
  bio = $10,
  birth_date = $9,
  last_location_long = $8,
  last_location_lat = $7,
  mobile = $6,
  last_active = $5,
  email = $4,
  sex = $3,
  interested_in = $2
WHERE id = $1
RETURNING *;