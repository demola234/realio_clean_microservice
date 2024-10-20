package usecase

import (
	"context"
	"errors"
	"fmt"
	interfaces "job_portal/authentication/interfaces/error"
	"job_portal/authentication/internal/domain/entity"
	"job_portal/authentication/internal/domain/repository"
	"job_portal/authentication/pkg/utils"
	"job_portal/authentication/pkg/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserUsecase defines the interface for user-related business logic.
type UserUsecase interface {
	RegisterUser(ctx context.Context, fullName, password, email string) (*entity.User, error)
	LoginUser(ctx context.Context, password, email string) (*entity.User, error)
	ChangePassword(ctx context.Context, currentPassword, newPassword, id string) error
	GenerateToken(ctx context.Context, email string) (string, error)
}

// userUsecase implements the UserUsecase interface.
type userUsecase struct {
	repo repository.UserRepository
}

// NewUserUsecase creates a new instance of userUsecase.
func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

// GenerateToken implements UserUsecase.
func (u *userUsecase) GenerateToken(ctx context.Context, email string) (string, error) {
	token, err := u.repo.CreateToken(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to generate token for email %s: %w", email, err)
	}
	return token, nil
}

// RegisterUser registers a new user.
func (u *userUsecase) RegisterUser(ctx context.Context, fullName string, password string, email string) (*entity.User, error) {
	// // Validate user input
	// violations := validateCreateUser(fullName, password, email)
	// if violations != nil {
	// 	return nil, interfaces.InvalidArgErr(violations)
	// }

	token, err := u.repo.CreateToken(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token for email %s: %w", email, err)
	}

	refreshToken, err := u.repo.CreateToken(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token for email %s: %w", email, err)
	}

	// Create a new user entity
	user := &entity.User{
		Email:        email,
		FullName:     fullName,
		Password:     password,
		AccessToken:  token,
		RefreshToken: refreshToken,
	}

	// Check if user already exists
	repoUser, _ := u.repo.GetUserByEmail(ctx, email)

	if repoUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Create the user
	err = u.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// LoginUser authenticates a user by email and password.
func (u *userUsecase) LoginUser(ctx context.Context, password string, email string) (*entity.User, error) {
	// Retrieve user by email
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user by email %s: %w", email, err)
	}

	// Check if the provided password matches the stored hash
	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials password incorrect %s: %w", email, err)
	}

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

	// Update the password in the repository
	err = u.repo.UpdatePassword(ctx, email, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
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
