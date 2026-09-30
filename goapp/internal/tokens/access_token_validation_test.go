package tokens

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

// The classic JWT bypass: re-sign the token with alg=none and hope the
// validator trusts the header. ParseAccessToken pins the method to HS256 in
// its keyfunc, which rejects this before any signature check happens.
func TestParseAccessToken_RejectsNoneAlgorithm(t *testing.T) {
	m, err := New(testSecret)
	require.NoError(t, err)

	unsigned := jwt.NewWithClaims(jwt.SigningMethodNone, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})

	tokenString, err := unsigned.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = m.ParseAccessToken(tokenString)
	require.Error(t, err)
}

// An asymmetric token must not be accepted either: without the method check,
// the HMAC secret would be handed to the RSA verifier as a public key.
func TestParseAccessToken_RejectsAsymmetricAlgorithm(t *testing.T) {
	m, err := New(testSecret)
	require.NoError(t, err)

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	rsaToken := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})

	tokenString, err := rsaToken.SignedString(key)
	require.NoError(t, err)

	_, err = m.ParseAccessToken(tokenString)
	require.Error(t, err)
}

// golang-jwt treats a missing exp as "not expired" rather than invalid, so a
// token with no exp claim would otherwise live forever. jwt.WithExpirationRequired
// closes that off; this holds it in place.
func TestParseAccessToken_RejectsTokenWithoutExpiry(t *testing.T) {
	m, err := New(testSecret)
	require.NoError(t, err)

	noExpiry := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  "user-123",
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := noExpiry.SignedString(m.secret)
	require.NoError(t, err)

	_, err = m.ParseAccessToken(tokenString)
	require.Error(t, err)
}

func TestParseAccessToken_RejectsTokenNotYetValid(t *testing.T) {
	m, err := New(testSecret)
	require.NoError(t, err)

	future := time.Now().Add(time.Hour)
	notYetValid := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			IssuedAt:  jwt.NewNumericDate(future),
			NotBefore: jwt.NewNumericDate(future),
			ExpiresAt: jwt.NewNumericDate(future.Add(time.Hour)),
		},
	})

	tokenString, err := notYetValid.SignedString(m.secret)
	require.NoError(t, err)

	_, err = m.ParseAccessToken(tokenString)
	require.Error(t, err)
}
