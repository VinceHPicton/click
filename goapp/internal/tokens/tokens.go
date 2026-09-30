package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultAccessTokenLifetime = 15 * time.Minute
	// clockSkewAllowance backdates NotBefore slightly so a token generated
	// on this instance isn't rejected by a validator whose clock is a
	// little behind.
	clockSkewAllowance = 1 * time.Minute
)

type Claims struct {
	jwt.RegisteredClaims
}

type Manager struct {
	secret  []byte
	ttl     time.Duration
	nowFunc func() time.Time
}

// Option configures a Manager at construction time.
type Option func(*Manager)

// WithTTL overrides the default access token lifetime.
func WithTTL(ttl time.Duration) Option {
	return func(m *Manager) {
		m.ttl = ttl
	}
}

// WithNowFunc overrides the clock the Manager uses to generate and validate
// tokens. Intended for tests (e.g. minting an already-expired token, or
// freezing time); production callers should not normally need this.
func WithNowFunc(now func() time.Time) Option {
	return func(m *Manager) {
		m.nowFunc = now
	}
}

func New(secret string, opts ...Option) (*Manager, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}

	m := &Manager{
		secret:  []byte(secret),
		ttl:     defaultAccessTokenLifetime,
		nowFunc: time.Now,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m, nil
}

func (m *Manager) GenerateAccessToken(userID string) (string, error) {
	now := m.nowFunc()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-clockSkewAllowance)),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(m.secret)
}

func (m *Manager) ParseAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return m.secret, nil
		},
		// Without this, golang-jwt treats a *missing* exp claim as valid
		// (not expired) rather than rejecting it. We always set ExpiresAt
		// ourselves, but this closes off that failure mode for good.
		jwt.WithExpirationRequired(),
		// Validate against the Manager's clock rather than the real
		// wall clock, so WithNowFunc actually controls expiry checks too.
		jwt.WithTimeFunc(m.nowFunc),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

// GenerateRefreshToken returns a cryptographically random, URL-safe token.
//
// NOTE: this now returns an error. crypto/rand.Read on the default reader
// essentially never fails on real operating systems, but silently ignoring
// a partial/failed read previously meant a broken token could be handed out
// instead of surfacing the failure. Callers must be updated.
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
