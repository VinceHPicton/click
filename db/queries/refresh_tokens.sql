-- name: CreateRefreshToken :one
INSERT INTO app.refresh_tokens (user_id, token_hash, expires_at, created_at)
VALUES (sqlc.arg(user_id), sqlc.arg(token_hash), NOW() + INTERVAL '30 days', NOW())
RETURNING *;

-- name: RotateRefreshToken :one
WITH revoked AS (
    UPDATE app.refresh_tokens
    SET revoked_at = NOW()
    WHERE refresh_tokens.user_id = sqlc.arg(user_id)
      AND refresh_tokens.token_hash = sqlc.arg(old_token_hash)
      AND refresh_tokens.revoked_at IS NULL
      AND refresh_tokens.expires_at > NOW()
    RETURNING *
),
inserted AS (
    INSERT INTO app.refresh_tokens (
        user_id,
        token_hash,
        expires_at,
        created_at
    )
    SELECT
        sqlc.arg(user_id),
        sqlc.arg(new_token_hash),
        NOW() + INTERVAL '30 days',
        NOW()
    FROM revoked
    RETURNING *
)
SELECT * FROM inserted;

-- name: GetValidTokenByHash :one
SELECT * FROM app.refresh_tokens
WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW();

-- name: GetTokenByHash :one
SELECT * FROM app.refresh_tokens
WHERE token_hash = $1;

-- name: GetRefreshTokens :many
SELECT * FROM app.refresh_tokens;