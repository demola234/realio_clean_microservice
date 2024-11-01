package token_maker

import (
	"testing"
	"time"

	"job_portal/shared/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewTokenMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomString(6)
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	userID := uuid.New().String()
	token, payload, err := maker.CreateToken(username, userID, duration)
	require.NoError(t, err)

	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	require.NoError(t, payload.Valid())

}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewTokenMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)
	email := utils.RandomString(6)

	userID := uuid.New().String()

	token, pasto_payload, err := maker.CreateToken(email, userID, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, pasto_payload)

}

func TestInvalidToken(t *testing.T) {
	maker, err := NewTokenMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	payload, err := maker.VerifyToken("invalid_token")
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Empty(t, payload)
}
