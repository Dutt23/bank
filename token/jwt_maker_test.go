package token

import (
	"github/dutt23/bank/util"
	"testing"
	"time"

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

	token, err := maker.CreateToken(username, duration)

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
