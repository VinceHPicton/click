package factory

import (
	"context"
	"vincehpicton/click/internal/db/sqlc"
)

func FakeAuthAttempt(ctx context.Context, q *sqlc.Queries, mob string, otc int32) (sqlc.AppAuthAttempt, error) {
	params := sqlc.AuthAttemptCreateWithOneTimeCodeParams{
		Mobile:      mob,
		OneTimeCode: otc,
	}

	return q.AuthAttemptCreateWithOneTimeCode(ctx, params)
}
