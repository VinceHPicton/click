package tokens

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testSecret = "super-secret-test-key"

func TestGenerateAndParseAccessToken(t *testing.T) {
	m, err := New(testSecret)
	require.NoError(t, err)
	const fakeUserID = "user-123"

	token, err := m.GenerateAccessToken(fakeUserID)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := m.ParseAccessToken(token)

	require.NoError(t, err)
	require.Equal(t, fakeUserID, claims.Subject)
}

func TestParseAccessToken_InvalidSignature(t *testing.T) {
	m1, err := New("secret-1")
	require.NoError(t, err)
	m2, err := New("secret-2")
	require.NoError(t, err)

	token, err := m1.GenerateAccessToken("user-123")
	require.NoError(t, err)

	_, err = m2.ParseAccessToken(token)

	require.Error(t, err)
}

func TestParseAccessToken_InvalidToken(t *testing.T) {
	m, err := New(testSecret)
	require.NoError(t, err)

	_, err = m.ParseAccessToken("this-is-not-a-jwt")

	require.Error(t, err)
}

func TestHashRefreshToken(t *testing.T) {
	hash1 := HashRefreshToken("abc")
	hash2 := HashRefreshToken("abc")
	hash3 := HashRefreshToken("xyz")

	require.Equal(t, hash1, hash2)
	require.NotEqual(t, hash1, hash3)
}

func TestExpiredToken(t *testing.T) {
	m, err := New(testSecret)
	require.NoError(t, err)
	m.ttl = -1 * time.Minute

	token, err := m.GenerateAccessToken("user-123")
	require.NoError(t, err)

	_, err = m.ParseAccessToken(token)

	require.Error(t, err)
}
