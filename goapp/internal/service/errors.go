package service

import "errors"

// Sentinel errors returned by the service for expected, caller-visible failures.
//
// Anything not in this list is an unexpected failure (a DB error, a constraint
// violation, a broken invariant) and should be treated by callers as an
// internal error: logged server-side, reported generically to the client.
//
// The service deliberately knows nothing about HTTP. Mapping these to status
// codes is the transport layer's job; see httpserver.writeError.
var (
	ErrInvalidAttempt = errors.New("auth attempt is invalid, expired, or already used")

	ErrInvalidCode = errors.New("incorrect one-time code")

	ErrUserNotFound = errors.New("no active user for mobile")

	ErrUserExists = errors.New("an active user already exists for mobile")

	// ErrInvalidRefreshToken covers unknown, revoked, expired, and
	// already-rotated refresh tokens. They are deliberately not distinguished:
	// telling a caller which one it was leaks whether a token ever existed.
	ErrInvalidRefreshToken = errors.New("refresh token is invalid, expired, or already used")

	ErrUserBanned = errors.New("user is banned")

	ErrUserDeleted = errors.New("user is deleted")
)
