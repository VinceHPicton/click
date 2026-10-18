package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) ConfirmLogin(ctx context.Context, attemptID uuid.UUID, oneTimeCode int32) (Session, error) {
	authAttempt, err := getAndConsumeAttempt(ctx, s.queries, attemptID, oneTimeCode)
	if err != nil {
		return Session{}, err
	}

	user, err := s.queries.GetActiveUserByMobile(ctx, authAttempt.Mobile)
	if err != nil {
		return Session{}, ErrUserNotFound
	}

	return issueSession(ctx, s.queries, s.tokenMgr, user.ID)
}
