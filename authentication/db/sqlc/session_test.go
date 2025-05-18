package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/demola234/authentication/pkg/utils"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T) Sessions {
	user := createRandomUser(t)

	arg := CreateSessionParams{
		SessionID:    uuid.New(),
		UserID:       user.ID, // Use the ID of the created user
		Token:        utils.RandomString(32),
		ExpiresAt:    time.Now().Add(24 * time.Hour).UTC(), // Set to UTC
		LastActivity: time.Now().UTC(),                     // Set to UTC
		IpAddress:    sql.NullString{String: "127.0.0.1", Valid: true},
		UserAgent:    sql.NullString{String: "Mozilla/5.0", Valid: true},
		Otp:          sql.NullString{String: "123456", Valid: true},
		OtpVerified:  sql.NullBool{Bool: false, Valid: true},
		OtpExpiresAt: sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
		OtpAttempts:  sql.NullInt32{Int32: 0, Valid: true},
		IsActive:     true,
		DeviceInfo:   pqtype.NullRawMessage{RawMessage: []byte(`{"device": "laptop"}`), Valid: true},
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, arg.SessionID, session.SessionID)
	require.Equal(t, user.ID, session.UserID)
	require.Equal(t, arg.Token, session.Token)
	require.True(t, session.IsActive)
	require.WithinDuration(t, arg.ExpiresAt, session.ExpiresAt.UTC(), time.Second)       // Compare as UTC
	require.WithinDuration(t, arg.LastActivity, session.LastActivity.UTC(), time.Second) // Compare as UTC

	return session
}

func TestCreateSession(t *testing.T) {
	session := createRandomSession(t)
	require.NotEmpty(t, session)
}

func TestGetSession(t *testing.T) {
	session := createRandomSession(t)

	session2, err := testQueries.GetSessionByID(context.Background(), session.SessionID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)
	require.Equal(t, session.SessionID, session2.SessionID)
	require.Equal(t, session.UserID, session2.UserID)
	require.Equal(t, session.Token, session2.Token)
	require.Equal(t, session.IsActive, session2.IsActive)
	require.WithinDuration(t, session.CreatedAt, session2.CreatedAt, time.Second)
	require.WithinDuration(t, session.ExpiresAt, session2.ExpiresAt, time.Second)
	require.WithinDuration(t, session.LastActivity, session2.LastActivity, time.Second)
}

func TestRevokeSession(t *testing.T) {
	session := createRandomSession(t)

	err := testQueries.RevokeSession(context.Background(), session.SessionID)
	require.NoError(t, err)

	// Verify the session is revoked
	session2, err := testQueries.GetSessionByID(context.Background(), session.SessionID)
	require.NoError(t, err)
	require.False(t, session2.IsActive)
	require.NotZero(t, session2.RevokedAt)
}

func TestUpdateSessionActivity(t *testing.T) {
	session := createRandomSession(t)

	newLastActivity := time.Now().Add(1 * time.Hour).UTC() // Set to UTC
	arg := UpdateSessionActivityParams{
		LastActivity: newLastActivity,
		IsActive:     session.IsActive,
		SessionID:    session.SessionID,
	}

	err := testQueries.UpdateSessionActivity(context.Background(), arg)
	require.NoError(t, err)

	// Retrieve and verify the session's last activity timestamp is updated
	session2, err := testQueries.GetSessionByID(context.Background(), session.SessionID)
	require.NoError(t, err)
	require.WithinDuration(t, newLastActivity, session2.LastActivity.UTC(), time.Second) // Compare as UTC
	require.Equal(t, session.IsActive, session2.IsActive)
}

func TestDeleteSession(t *testing.T) {
	session := createRandomSession(t)

	err := testQueries.DeleteSession(context.Background(), session.SessionID)
	require.NoError(t, err)

	// Verify the session no longer exists
	session2, err := testQueries.GetSessionByID(context.Background(), session.SessionID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, session2)
}

func TestUpdateOtp(t *testing.T) {
	session := createRandomSession(t)
	newOtp := "789012"

	arg := UpdateSessionParams{
		UserID:       session.UserID,
		Otp:          sql.NullString{String: newOtp, Valid: true},
		OtpVerified:  sql.NullBool{Bool: true, Valid: true},
		OtpExpiresAt: sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
		OtpAttempts:  sql.NullInt32{Int32: 0, Valid: true},
	}

	_, err := testQueries.UpdateSession(context.Background(), arg)

	require.NoError(t, err)
	session2, err := testQueries.GetSessionByID(context.Background(), session.SessionID)
	require.NoError(t, err)
	require.Equal(t, newOtp, session2.Otp.String)
	require.Equal(t, true, session2.OtpVerified.Bool)
	require.Equal(t, int32(0), session2.OtpAttempts.Int32)
	require.Equal(t, session.SessionID, session2.SessionID)
	require.Equal(t, session.UserID, session2.UserID)
}
