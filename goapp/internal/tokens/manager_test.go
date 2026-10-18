package tokens

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew_RejectsEmptySecret(t *testing.T) {
	m, err := New("")
	require.Error(t, err)
	require.Nil(t, m)
}

func TestWithTTL_ControlsExpiry(t *testing.T) {
	const ttl = 3 * time.Hour
	issuedAt := time.Date(2026, time.September, 16, 12, 0, 0, 0, time.UTC)

	m, err := New(testSecret,
		WithTTL(ttl),
		WithNowFunc(func() time.Time { return issuedAt }),
	)
	require.NoError(t, err)

	token, err := m.GenerateAccessToken("user-123")
	require.NoError(t, err)

	claims, err := m.ParseAccessToken(token)
	require.NoError(t, err)
	require.Equal(t, issuedAt.Add(ttl).Unix(), claims.ExpiresAt.Unix())
}

func TestWithTTL_NegativeTTLMintsExpiredToken(t *testing.T) {
	m, err := New(testSecret, WithTTL(-1*time.Minute))
	require.NoError(t, err)

	token, err := m.GenerateAccessToken("user-123")
	require.NoError(t, err)

	_, err = m.ParseAccessToken(token)
	require.Error(t, err)
}

// clockSkewAllowance backdates NotBefore so a freshly minted token is not
// rejected by a validator running slightly behind. Without it a validator one
// second in the past would reject every token it was handed.
func TestGenerateAccessToken_TolerantOfValidatorClockBehind(t *testing.T) {
	issuedAtYear, issuedAtMonth, issuedAtDay := time.Now().Date()
	issuedAt := time.Date(issuedAtYear, issuedAtMonth, issuedAtDay, 12, 0, 0, 0, time.UTC)

	issuer, err := New(testSecret, WithNowFunc(func() time.Time { return issuedAt }))
	require.NoError(t, err)

	token, err := issuer.GenerateAccessToken("user-123")
	require.NoError(t, err)

	validator, err := New(testSecret, WithNowFunc(func() time.Time {
		return issuedAt.Add(-30 * time.Second)
	}))
	require.NoError(t, err)

	_, err = validator.ParseAccessToken(token)
	require.NoError(t, err)
}
