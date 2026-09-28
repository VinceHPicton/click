package service

import (
	"context"
	"vincehpicton/click/internal/db/sqlc"
)

func (s *Service) StartLogin(ctx context.Context, mobile string) (sqlc.AppAuthAttempt, error) {
	attempt, err := s.queries.AuthAttemptCreate(ctx, mobile)
	if err != nil {
		return sqlc.AppAuthAttempt{}, err
	}

	return attempt, nil
}
