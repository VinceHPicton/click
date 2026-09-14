package service

import "errors"

var (
	ErrInvalidAttempt = errors.New("auth attempt is invalid, expired, or already used")

	ErrInvalidCode = errors.New("incorrect one-time code")

	ErrUserNotFound = errors.New("no active user for mobile")

	ErrUserExists = errors.New("an active user already exists for mobile")
)
