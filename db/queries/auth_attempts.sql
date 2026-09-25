-- name: AuthAttemptCreate :one

INSERT INTO app.auth_attempts (
  mobile
) VALUES (
  $1
)
RETURNING *;

-- name: AuthAttemptCreateWithOneTimeCode :one
INSERT INTO app.auth_attempts (
  mobile,
  one_time_code
) VALUES (
  $1,
  $2
)
RETURNING *;

-- name: GetAuthAttempts :many
SELECT * FROM app.auth_attempts;

-- name: GetAuthAttempt :one
SELECT * FROM app.auth_attempts
WHERE id = $1;

-- name: GetValidAuthAttempt :one
SELECT * FROM app.auth_attempts
WHERE id = $1
AND used_at IS NULL
AND created_at >= NOW() - INTERVAL '2 minutes';

-- name: ConsumeAuthAttempt :exec
UPDATE app.auth_attempts
SET used_at = NOW()
WHERE id = sqlc.arg(id)
AND used_at IS NULL;

-- name: ConsumeAuthAttemptByID :exec
UPDATE app.auth_attempts
SET used_at = NOW()
WHERE id = $1
  AND used_at IS NULL
  AND created_at >= NOW() - INTERVAL '2 minutes';

-- Test support.
-- name: SetAuthAttemptCreatedAt :exec
UPDATE app.auth_attempts
SET created_at = sqlc.arg(created_at)
WHERE id = sqlc.arg(id);

-- Test support. 
-- name: SetAuthAttemptUsedAt :exec
UPDATE app.auth_attempts
SET used_at = sqlc.arg(used_at)
WHERE id = sqlc.arg(id);