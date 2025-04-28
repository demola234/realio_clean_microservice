package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	mockdb "job_portal/authentication/db/mock"
	db "job_portal/authentication/db/sqlc"
	"job_portal/authentication/internal/domain/entity"
	"job_portal/authentication/pkg/utils"
)

func TestGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	repo := NewUserRepository(store)

	user := db.Users{
		ID:        uuid.New(),
		Email:     utils.RandomEmail(),
		Name:      utils.RandomOwner(),
		Password:  sql.NullString{String: utils.RandomString(10), Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Role:      sql.NullString{String: utils.RandomRole(), Valid: true},
		Phone:     sql.NullString{String: utils.RandomPhoneNumber(), Valid: true},
	}

	store.EXPECT().
		GetUser(gomock.Any(), user.Email).
		Times(1).
		Return(user, nil)

	result, err := repo.GetUserByEmail(context.Background(), user.Email)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, user.ID, result.ID)
	require.Equal(t, user.Email, result.Email)
	require.Equal(t, user.Name, result.FullName)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	repo := NewUserRepository(store)

	user := &entity.User{
		ID:       uuid.New(),
		Email:    utils.RandomEmail(),
		FullName: utils.RandomOwner(),
		Password: utils.RandomString(10),
		Role:     utils.RandomRole(),
		Phone:    utils.RandomPhoneNumber(),
	}

	store.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Times(1).
		Return(db.Users{}, nil)

	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
}

func TestCreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	repo := NewUserRepository(store)

	email := utils.RandomEmail()
	userID := uuid.New().String()

	token, err := repo.CreateToken(context.Background(), email, userID)

	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestUpdatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	repo := NewUserRepository(store)

	userID := uuid.New().String()
	newPassword := utils.RandomString(10)

	store.EXPECT().
		UpdateUser(gomock.Any(), gomock.Any()).
		Times(1).
		Return(db.Users{}, nil)

	err := repo.UpdatePassword(context.Background(), userID, newPassword)
	require.NoError(t, err)
}
