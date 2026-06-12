-- name: AuthAttemptCreate :one

INSERT INTO app.auth_attempts (
  mobile
) VALUES (
  $1
)
RETURNING *;

-- name: GetAuthAttempts :many

SELECT * FROM app.auth_attempts;

-- Note: you can use :exec if it doesnt return anything
-- name: AuthAttemptCreateUser :one

WITH consumed_attempt AS (
  
    UPDATE app.auth_attempts AS ra
    SET used_at = NOW()
    WHERE ra.id = sqlc.arg(id)
      AND ra.one_time_code = sqlc.arg(one_time_code)
      AND ra.used_at IS NULL
      AND ra.created_at >= NOW() - INTERVAL '2 minutes'
    RETURNING id, mobile
),
new_user AS (
    INSERT INTO app.users (mobile)
    SELECT mobile
    FROM consumed_attempt
    RETURNING id
)
SELECT id AS user_id
FROM new_user;

-- name: GetAuthAttempt :one