package usecase

import (
	"context"
	"job_portal/authentication/internal/domain/entity"
	"job_portal/authentication/pkg/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

// GetOtp implements repository.UserRepository.
func (m *MockUserRepository) GetOtp(ctx context.Context, id string) (*entity.UpdateOtp, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UpdateOtp), args.Error(1)
}

// UpdateSession implements repository.UserRepository.
func (m *MockUserRepository) UpdateSession(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) CreateToken(ctx context.Context, email string, userID string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, email string, hashedPassword string) error {
	args := m.Called(ctx, email, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserSession(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockUserRepository) DeleteSession(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateOtp(ctx context.Context, updateOtp *entity.UpdateOtp) error {
	args := m.Called(ctx, updateOtp)
	return args.Error(0)
}

func TestRegisterUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	fullName := "Test User"
	password := "password123"
	email := "test@example.com"
	role := "user"
	phone := "1234567890"

	// Mock behavior
	mockRepo.On("GetUserByEmail", ctx, email).Return(nil, nil)
	mockRepo.On("CreateToken", ctx, email).Return("test-token", nil)
	mockRepo.On("CreateSession", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)
	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	// Execute test
	user, session, err := useCase.RegisterUser(ctx, fullName, password, email, role, phone)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, session)
	require.Equal(t, email, user.Email)
	require.Equal(t, fullName, user.FullName)
	require.Equal(t, role, user.Role)
	require.Equal(t, phone, user.Phone)
}

func TestLoginUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	password := "password123"
	hashedPassword, _ := utils.HashPassword(password)
	email := "test@example.com"

	mockUser := &entity.User{
		ID:       uuid.New(),
		Email:    email,
		Password: hashedPassword,
		FullName: "Test User",
	}

	// Mock behavior
	mockRepo.On("GetUserByEmail", ctx, email).Return(mockUser, nil)

	// Execute test
	user, err := useCase.LoginUser(ctx, password, email)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, email, user.Email)
}

func TestChangePassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	currentPassword := "oldpassword123"
	newPassword := "newpassword123"
	email := "test@example.com"
	hashedOldPassword, _ := utils.HashPassword(currentPassword)

	mockUser := &entity.User{
		ID:       uuid.New(),
		Email:    email,
		Password: hashedOldPassword,
	}

	// Mock behavior
	mockRepo.On("GetUserByEmail", ctx, email).Return(mockUser, nil)
	mockRepo.On("UpdatePassword", ctx, email, mock.AnythingOfType("string")).Return(nil)

	// Execute test
	err := useCase.ChangePassword(ctx, currentPassword, newPassword, email)

	// Assertions
	require.NoError(t, err)
}

func TestGetSession(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	sessionID := uuid.New().String()
	mockSession := &entity.Session{
		SessionID: uuid.MustParse(sessionID),
		UserID:    uuid.New(),
		Token:     "test-token",
		IsActive:  true,
	}

	// Mock behavior
	mockRepo.On("GetUserSession", ctx, sessionID).Return(mockSession, nil)

	// Execute test
	session, err := useCase.GetSession(ctx, sessionID)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, session)
	require.Equal(t, mockSession.SessionID, session.SessionID)
	require.Equal(t, mockSession.Token, session.Token)
}

func TestLogOut(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	userID := uuid.New().String()

	// Mock behavior
	mockRepo.On("DeleteSession", ctx, userID).Return(nil)

	// Execute test
	err := useCase.LogOut(ctx, userID)

	// Assertions
	require.NoError(t, err)
}

func TestResendOtp(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	email := "test@example.com"

	// Mock behavior
	mockRepo.On("UpdateOtp", ctx, mock.AnythingOfType("*entity.UpdateOtp")).Return(nil)

	// Execute test
	err := useCase.ResendOtp(ctx, email)

	// Assertions
	require.NoError(t, err)
}

// func TestVerifyOtp(t *testing.T) {
// 	mockRepo := new(MockUserRepository)
// 	useCase := NewUserUsecase(mockRepo)
// 	ctx := context.Background()

// 	email := "test@example.com"
// 	otp := "123456"
// 	hashedOtp, _ := utils.HashPassword(otp)
// 	mockUser := &entity.User{
// 		ID:       uuid.New(),
// 		Email:    email,
// 		Password: hashedOtp,
// 	}

// }
