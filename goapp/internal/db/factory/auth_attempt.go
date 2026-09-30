package factory

import (
	"context"
	"time"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/service"
)

func FakeAuthAttempt(ctx context.Context, q *sqlc.Queries, mob string, otc int32) (sqlc.AppAuthAttempt, error) {
	params := sqlc.AuthAttemptCreateWithOneTimeCodeParams{
		Mobile:      mob,
		OneTimeCode: otc,
	}

	return q.AuthAttemptCreateWithOneTimeCode(ctx, params)
}

func FakeAuthAttemptWithServiceExpiry(ctx context.Context, q *sqlc.Queries, mob string) (sqlc.AppAuthAttempt, error) {
	return q.AuthAttemptCreate(ctx, sqlc.AuthAttemptCreateParams{
		Mobile:    mob,
		ExpiresAt: time.Now().Add(service.AuthAttemptValidityWindowMinutes * time.Minute),
	})
}
