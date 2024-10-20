package repository

import (
	"context"
	"job_portal/authentication/internal/domain/entity"
)

// AuthRepository defines the repository contract for user-related operations.
type UserRepository interface {
	// GetUserByEmail fetches a user by their email address.
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// CreateUser saves a new user to the data store.
	CreateUser(ctx context.Context, user *entity.User) error

	// UpdatePassword updates the password for an existing user.
	UpdatePassword(ctx context.Context, email string, newPassword string) error

	CreateToken(ctx context.Context, email string) (string, error)
}
