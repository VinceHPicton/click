-- name: RegisterAttemptCreate :one
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO register_attempts (
  mobile
) VALUES (
  $1
)
RETURNING *;