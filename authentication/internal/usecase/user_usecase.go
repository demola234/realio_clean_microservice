package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	interfaces "github.com/demola234/authentication/infrastructure/error"
	"github.com/demola234/authentication/internal/domain/entity"
	"github.com/demola234/authentication/internal/domain/repository"
	"github.com/demola234/authentication/pkg/utils"
	"github.com/demola234/authentication/pkg/val"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserUsecase defines the interface for user-related business logic.
type UserUsecase interface {
	RegisterUser(ctx context.Context, fullName string, password string, email string, role string, phone string) (*entity.User, *entity.Session, error)
	RegisterWithOAuth(ctx context.Context, provider, token string) (*entity.User, *entity.Session, error)
	LoginUser(ctx context.Context, password, email string) (*entity.User, error)
	ChangePassword(ctx context.Context, currentPassword, newPassword, id string) error
	GetSession(ctx context.Context, id string) (*entity.Session, error)
	GenerateToken(ctx context.Context, email string, userID string) (string, error)
	ResendOtp(ctx context.Context, email string) error
	GetUser(ctx context.Context, userId string) (*entity.User, error)
	LogOut(ctx context.Context, userId string) error
	VerifyOtp(ctx context.Context, email string, otp string) (bool, error)
}

// userUsecase implements the UserUsecase interface.
type userUsecase struct {
	repo repository.UserRepository
}


// NewUserUsecase creates a new instance of userUsecase.
func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}


// RegisterWithOAuth implements UserUsecase.
func (u *userUsecase) RegisterWithOAuth(ctx context.Context, provider string, token string) (*entity.User, *entity.Session, error) {
	panic("unimplemented")
}


// GenerateToken implements UserUsecase.
func (u *userUsecase) GenerateToken(ctx context.Context, email string, userID string) (string, error) {
	token, err := u.repo.CreateToken(ctx, email, userID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token for email %s: %w", email, err)
	}
	return token, nil
}

// RegisterUser registers a new user.
func (u *userUsecase) RegisterUser(ctx context.Context, fullName, password, email, role, phone string) (*entity.User, *entity.Session, error) {
	// Validate user input (uncomment if needed)
	// violations := validateCreateUser(fullName, password, email)
	// if violations != nil {
	// 	return nil, nil, interfaces.InvalidArgErr(violations)
	// }

	// Check if user with the same email already exists
	existingUser, err := u.repo.GetUserByEmail(ctx, email)
	userID := uuid.New()
	if err == nil && existingUser != nil {
		return nil, nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Generate a new token for the session
	token, err := u.repo.CreateToken(ctx, email, userID.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate token for email %s: %w", email, err)
	}

	// Create a new user entity with a unique ID
	user := &entity.User{
		ID:       userID,
		Email:    email,
		FullName: fullName,
		Password: password,
		Role:     role,
		Phone:    phone,
	}

	// Gather metadata for session
	metaData := utils.ExtractMetaData(ctx)

	// Create a new session entity with a unique ID and expiration time
	session := &entity.Session{
		SessionID:    uuid.New(),
		UserID:       user.ID,
		Token:        token,
		CreatedAt:    time.Now().UTC(),
		ExpiresAt:    time.Now().Add(24 * time.Hour).UTC(), // Set to UTC
		LastActivity: time.Now().UTC(),
		IpAddress:    metaData.ClientIP,
		UserAgent:    metaData.UserAgent,
		IsActive:     true,
		Otp:          utils.RandomOtp(),
		OTPVerified:  false,
		OtpExpiresAt: time.Now().Add(5 * time.Minute),
		OtpAttempts:  0,
	}

	// Save user in the repository
	if err := u.repo.CreateUser(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Save session in the repository
	if err := u.repo.CreateSession(ctx, session); err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	return user, session, nil
}

// LoginUser authenticates a user by email and password.
func (u *userUsecase) LoginUser(ctx context.Context, password string, email string) (*entity.User, error) {
	// Retrieve user by email
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	// Check if the provided password matches the stored hash
	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		return nil, err
	}

	// session, err := u.repo.GetUserSession(ctx, user.ID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to retrieve user session: %w", err)
	// }

	// // Check if the user is active
	// if !session.IsActive {
	// 	return nil, fmt.Errorf("user is not active")
	// }

	// Check if user is verified
	// if !session.OTPVerified {
	// 	return nil, fmt.Errorf("user is not verified: %+v", session)
	// }

	return user, nil
}

// ChangePassword updates a user's password.
func (u *userUsecase) ChangePassword(ctx context.Context, currentPassword string, newPassword string, email string) error {
	// Retrieve the user by ID to verify the current password
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to retrieve user by ID %s: %w", email, err)
	}

	// Check if current password is correct
	if err := utils.CheckPassword(currentPassword, user.Password); err != nil {
		return fmt.Errorf("current password is incorrect: %w", err)
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	session, err := u.repo.GetUserSession(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve user session: %w", err)
	}

	// Check if the user is active
	if !session.IsActive {
		return fmt.Errorf("user is not active")
	}

	// Check if user is verified
	if !session.OTPVerified {
		return fmt.Errorf("user is not verified")
	}

	// Update the password in the repository
	err = u.repo.UpdatePassword(ctx, email, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// GetUser implements UserUsecase.
func (u *userUsecase) GetUser(ctx context.Context, userId string) (*entity.User, error) {
	// Retrieve user by email
	user, err := u.repo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user by userId %s: %w", userId, err)
	}

	return user, nil
}

// LogOut implements UserUsecase.
func (u *userUsecase) LogOut(ctx context.Context, userId string) error {
	// Retrieve user by email
	err := u.repo.DeleteSession(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to retrieve user by userId %s: %w", userId, err)
	}

	return nil
}

// ResendOtp implements UserUsecase.
func (u *userUsecase) ResendOtp(ctx context.Context, email string) error {
	// Retrieve user by email
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to retrieve user by email %s: %w", email, err)
	}

	//

	// Retrieve user by email
	err = u.repo.UpdateOtp(ctx, &entity.UpdateOtp{
		Otp:          utils.RandomOtp(),
		Email:        user.Email,
		OtpExpiresAt: time.Now().Add(time.Minute * 10),
	})

	if err != nil {
		return fmt.Errorf("failed to update otp: %w", err)
	}

	return nil
}

// GetSession implements UserUsecase.
func (u *userUsecase) GetSession(ctx context.Context, id string) (*entity.Session, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID: %w", err)
	}
	// Retrieve user by email
	session, err := u.repo.GetUserSession(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user by userId %s: %w", id, err)
	}

	return session, nil

}

// VerifyOtp implements UserUsecase.
func (u *userUsecase) VerifyOtp(ctx context.Context, email string, otp string) (bool, error) {
	// Retrieve user by email
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user by email %s: %w", email, err)
	}

	session, err := u.repo.GetUserSession(ctx, user.ID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user session: %w", err)
	}

	otpUpdate, err := u.repo.GetOtp(ctx, user.ID.String())
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user session: %w", err)
	}

	if otpUpdate.Otp != otp {
		return false, fmt.Errorf("invalid otp %s: %s", otp, otpUpdate.Otp)
	}

	if session.OtpExpiresAt.After(time.Now()) {
		return false, fmt.Errorf("otp expired")
	}

	if session.OtpAttempts >= 10 {
		return false, fmt.Errorf("otp attempts exceeded")
	}

	// Update the password in the repository
	err = u.repo.UpdateOtp(ctx, &entity.UpdateOtp{
		OtpAttempts: session.OtpAttempts + 1,
		Email:       user.Email,
		OTPVerified: true,
	})

	if err != nil {
		return false, fmt.Errorf("failed to update password: %w", err)
	}

	return true, nil

}

// validateCreateUser checks if the provided inputs are valid.
func validateCreateUser(fullName string, password string, email string) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(email); err != nil {
		violations = append(violations, interfaces.FieldViolation("email", fmt.Errorf("invalid email format: %w", err)))
	}

	if err := val.ValidateFullName(fullName); err != nil {
		violations = append(violations, interfaces.FieldViolation("full_name", fmt.Errorf("invalid full name: %w", err)))
	}

	if err := val.ValidatePassword(password); err != nil {
		violations = append(violations, interfaces.FieldViolation("password", fmt.Errorf("invalid password: %w", err)))
	}

	return violations
}
