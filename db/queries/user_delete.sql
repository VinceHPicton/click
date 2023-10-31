-- name: UserDelete :exec
DELETE FROM users WHERE id = $1;