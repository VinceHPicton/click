package tokens

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRefreshToken_Unique(t *testing.T) {
	const draws = 100

	seen := make(map[string]struct{}, draws)
	for range draws {
		token, err := GenerateRefreshToken()
		require.NoError(t, err)
		_, duplicate := seen[token]
		require.False(t, duplicate, "refresh token collision")
		seen[token] = struct{}{}
	}
}

func TestGenerateRefreshToken_Carries32BytesOfEntropy(t *testing.T) {
	token, err := GenerateRefreshToken()
	require.NoError(t, err)

	raw, err := base64.RawURLEncoding.DecodeString(token)
	require.NoError(t, err)
	require.Len(t, raw, 32)
}

// The token travels in JSON bodies and is compared verbatim, so it must not
// need escaping: no padding, and no characters outside the URL-safe alphabet.
func TestGenerateRefreshToken_URLSafeAndUnpadded(t *testing.T) {
	token, err := GenerateRefreshToken()
	require.NoError(t, err)

	require.NotContains(t, token, "=")
	require.NotContains(t, token, "+")
	require.NotContains(t, token, "/")

	const urlSafeAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	for _, r := range token {
		require.True(t, strings.ContainsRune(urlSafeAlphabet, r), "unexpected character %q", r)
	}
}

func TestHashRefreshToken_IsHexSHA256(t *testing.T) {
	hash := HashRefreshToken("abc")
	require.Len(t, hash, 64)
	require.NotContains(t, hash, "abc")
}
