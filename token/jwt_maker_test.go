package token

import (
	"errors"
	"github/dutt23/bank/util"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {

	maker, err := newJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	// 1 minute
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, _, err := maker.CreateToken(username, duration)

	require.NotEmpty(t, token)
	require.NoError(t, err)

	payload, err := maker.Validate(token)

	require.NotEmpty(t, payload)
	require.NoError(t, err)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestJWTExpiredToken(t *testing.T) {
	maker, err := newJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	token, _, err := maker.CreateToken(username, -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.Validate(token)

	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, errors.New("token has expired").Error())
}

func TestIvalidJWTToken(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwttoken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwttoken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	maker, err := newJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	validatedToken, tokenErr := maker.Validate(token)
	require.Error(t, tokenErr)
	require.Nil(t, validatedToken)
	require.EqualError(t, tokenErr, errors.New("signing method is not valid here").Error())
}
