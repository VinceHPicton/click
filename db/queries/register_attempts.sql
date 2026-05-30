-- name: RegisterAttemptCreate :one

INSERT INTO register_attempts (
  mobile
) VALUES (
  $1
)
RETURNING *;

-- Note: you can use :exec if it doesnt return anything
-- name: RegisterAttemptConfirm :one

UPDATE register_attempts
SET
    used_at = NOW()
WHERE id = sqlc.arg(id)
  AND one_time_code = sqlc.arg(one_time_code)
  AND used_at IS NULL
  AND created_at >= NOW() - INTERVAL '2 minutes'
RETURNING id;