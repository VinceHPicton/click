package service

import (
	"context"
	"time"
	"vincehpicton/click/internal/db/sqlc"
)

func (s *Service) StartLogin(ctx context.Context, mobile string) (sqlc.AppAuthAttempt, error) {
	attempt, err := s.queries.AuthAttemptCreate(ctx, sqlc.AuthAttemptCreateParams{
		Mobile:    mobile,
		ExpiresAt: time.Now().Add(AuthAttemptValidityWindowMinutes * time.Minute),
	})
	if err != nil {
		return sqlc.AppAuthAttempt{}, err
	}

	return attempt, nil
}
