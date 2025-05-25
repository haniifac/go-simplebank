package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/haniifac/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	t.Logf("token: %s", token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	t.Logf("payload: %v", payload)

	// Check payload values
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt.Time, time.Second)
}

func TestPasetoTokenExpired(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Contains(t, err.Error(), jwt.ErrTokenExpired.Error())
	require.Nil(t, payload)
}

func TestInvalidKeyPasetoToken(t *testing.T) {
	_, err := NewPasetoMaker(util.RandomString(31))
	require.EqualError(t, err, ErrInvalidPasetoKeySize.Error())
}

func TestVerifyToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	// Test invalid token
	invalidToken := "v2.invalid.token.here"
	payload, err := maker.VerifyToken(invalidToken)
	require.Error(t, err)
	require.Nil(t, payload)
}
