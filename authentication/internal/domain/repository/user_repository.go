package repository

import (
	"context"
	"job_portal/authentication/internal/domain/entity"

	"github.com/google/uuid"
)

// AuthRepository defines the repository contract for user-related operations.
type UserRepository interface {
	// GetUserByEmail fetches a user by their email address.
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// CreateUser saves a new user to the data store.
	CreateUser(ctx context.Context, user *entity.User) error

	// UpdatePassword updates the password for an existing user.
	UpdatePassword(ctx context.Context, email string, newPassword string) error

	// CreateToken generates a new access token for a user.
	CreateToken(ctx context.Context, email string, userID string) (string, error)

	// GetUserByID retrieves a user by their ID.
	GetUserByID(ctx context.Context, id string) (*entity.User, error)

	// GetUserSession retrieves a user session by its ID.
	GetUserSession(ctx context.Context, id uuid.UUID) (*entity.Session, error)

	// CreateSession creates a new user session.
	CreateSession(ctx context.Context, session *entity.Session) error

	// UpdateSession updates an existing user session.
	UpdateSession(ctx context.Context, session *entity.Session) error

	// DeleteSession deletes a user session by its ID.
	DeleteSession(ctx context.Context, id string) error

	// UpdateOtp updates the OTP for a user session.
	UpdateOtp(ctx context.Context, session *entity.UpdateOtp) error

	// GetOtp retrieves the OTP for a user session.
	GetOtp(ctx context.Context, id string) (*entity.UpdateOtp, error)
}
