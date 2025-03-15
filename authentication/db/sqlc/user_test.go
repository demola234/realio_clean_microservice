package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"job_portal/authentication/pkg/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) Users {
	hashPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		ID:             uuid.New(),
		Name:           utils.RandomOwner(),
		Email:          utils.RandomEmail(),
		Password:       hashPassword,
		ProfilePicture: sql.NullString{String: utils.RandomProfilePicture(), Valid: true},
		Username:       utils.RandomOwner(),
		Bio:            sql.NullString{String: utils.RandomBio(), Valid: true},
		Role:           sql.NullString{String: utils.RandomRole(), Valid: true},
		Phone:          sql.NullString{String: utils.RandomPhoneNumber(), Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.Role, user.Role)
	require.NotZero(t, user.CreatedAt)

	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.Email)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Name, user2.Name)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Bio, user2.Bio)
	require.Equal(t, user.Phone, user2.Phone)
	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.ProfilePicture, user2.ProfilePicture)
	require.Equal(t, user.Password, user2.Password)
	require.Equal(t, user.Role, user2.Role)
	require.WithinDuration(t, user.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, user.UpdatedAt.Time, user2.UpdatedAt.Time, time.Second)
}

func TestChangePassword(t *testing.T) {
	user := createRandomUser(t)
	newPassword := utils.RandomString(8)
	arg := ChangePasswordParams{
		ID:       user.ID,
		Password: newPassword,
	}

	updatedUser, err := testQueries.ChangePassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, newPassword, updatedUser.Password)
	require.Equal(t, user.Role, updatedUser.Role)
	require.WithinDuration(t, user.CreatedAt.Time, updatedUser.CreatedAt.Time, time.Second)
	require.WithinDuration(t, time.Now(), updatedUser.UpdatedAt.Time, time.Second)
}

func TestCheckEmailExists(t *testing.T) {
	user := createRandomUser(t)
	exists, err := testQueries.CheckEmailExists(context.Background(), user.Email)

	require.NoError(t, err)
	require.True(t, exists)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)
	newEmail := utils.RandomEmail()
	arg := UpdateUserParams{
		Name:  user.Name,
		Email: newEmail,
		Role:  sql.NullString{String: utils.RandomRole(), Valid: true},
		Phone: sql.NullString{String: utils.RandomPhoneNumber(), Valid: true},
		ID:    user.ID,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.ID, updatedUser.ID)
	require.Equal(t, arg.Name, updatedUser.Name)
	require.Equal(t, arg.Email, updatedUser.Email)
	require.Equal(t, arg.Role, updatedUser.Role)
	require.WithinDuration(t, user.CreatedAt.Time, updatedUser.CreatedAt.Time, time.Second)
	require.WithinDuration(t, time.Now(), updatedUser.UpdatedAt.Time, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	// Verify user no longer exists
	user2, err := testQueries.GetUser(context.Background(), user.Email)
	require.Error(t, err)
	require.Empty(t, user2)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
