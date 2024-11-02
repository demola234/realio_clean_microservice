package repository

import (
	"context"
	"database/sql"
	"fmt"
	"job_portal/authentication/config"
	db "job_portal/authentication/db/sqlc" // SQLC generated code for interacting with the database
	token "job_portal/api_gateway/interfaces/middleware/token_maker"
	"job_portal/authentication/internal/domain/entity"
	"job_portal/authentication/pkg/utils"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// UserRepository implements the AuthRepository interface.
// This struct interacts with the database using SQLC-generated code.
type UserRepository struct {
	store db.Store
}

// CreateToken implements repository.UserRepository.
func (r *UserRepository) CreateToken(ctx context.Context, email string, userID string) (string, error) {
	// Load configuration
	configs, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}

	tokenMaker, err := token.NewTokenMaker(configs.TokenSymmetricKey)

	if err != nil {
		log.Fatalf("Failed to load env file: %s", err)
	}

	duration := time.Hour * 24
	accessToken, _, err := tokenMaker.CreateToken(email, userID, duration)
	if err != nil {
		return "", fmt.Errorf("some went wrong: %d", err)
	}

	return accessToken, nil
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(store db.Store) *UserRepository {
	return &UserRepository{
		store: store,
	}
}

// GetUserByEmail retrieves a user by their email from the database.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {

	userDetails, err := r.store.GetUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user by email %s: %w", email, err)
	}

	return &entity.User{
		ID:        userDetails.ID,
		FullName:  userDetails.Name,
		Email:     userDetails.Email,
		CreatedAt: userDetails.CreatedAt.Time,
		Password:  userDetails.Password,
		Role:      userDetails.Role.String,
		Phone:     userDetails.Phone.String,
		UpdatedAt: userDetails.UpdatedAt.Time,
	}, nil
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = r.store.CreateUser(ctx, db.CreateUserParams{
		Email:    user.Email,
		Name:     user.FullName,
		Password: hashedPassword,
		ID:       user.ID,
		Role:     sql.NullString{String: user.Role, Valid: true},
		Phone:    sql.NullString{String: user.Phone, Valid: true},
	})

	if err != nil {
		return err
	}

	return nil
}

// UpdatePassword updates a user's password in the database.
func (r *UserRepository) UpdatePassword(ctx context.Context, userID string, newPassword string) error {

	_, err := r.store.UpdateUser(ctx, db.UpdateUserParams{
		Password: newPassword,
	})

	return err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	userDetails, err := r.store.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:        userDetails.ID,
		FullName:  userDetails.Name,
		Email:     userDetails.Email,
		CreatedAt: userDetails.CreatedAt.Time,
		Password:  userDetails.Password,
		Role:      userDetails.Role.String,
		Phone:     userDetails.Phone.String,
		UpdatedAt: userDetails.UpdatedAt.Time,
	}, nil

}

func (r *UserRepository) GetUserSession(ctx context.Context, sessionID uuid.UUID) (*entity.Session, error) {

	// Retrieve session details from the store
	sessionDetails, err := r.store.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// // Convert sql.NullTime to *time.Time
	// var revokedAt *time.Time
	// if sessionDetails.RevokedAt.Valid {
	// 	revokedAt = &sessionDetails.RevokedAt.Time

	// }

	// Map session details to the entity.Session struct
	return &entity.Session{
		SessionID:    sessionDetails.SessionID,
		UserID:       sessionDetails.UserID,
		Token:        sessionDetails.Token,
		LastActivity: sessionDetails.LastActivity,
		ExpiresAt:    sessionDetails.ExpiresAt,
		CreatedAt:    sessionDetails.CreatedAt,
		DeviceInfo:   &sessionDetails.DeviceInfo,
		UserAgent:    sessionDetails.UserAgent.String,
		IpAddress:    sessionDetails.IpAddress.String,
		IsActive:     sessionDetails.IsActive,
		// RevokedAt:    revokedAt,
	}, nil
}

func (r *UserRepository) CreateSession(ctx context.Context, session *entity.Session) error {

	// Convert RevokedAt to sql.NullTime
	var revokedAt sql.NullTime
	if session.RevokedAt != nil {
		revokedAt = sql.NullTime{Time: *session.RevokedAt, Valid: true}
	} else {
		revokedAt = sql.NullTime{Valid: false}
	}

	// Call CreateSession with the mapped parameters
	_, err := r.store.CreateSession(ctx, db.CreateSessionParams{
		SessionID:    session.SessionID,
		UserID:       session.UserID,
		Token:        session.Token,
		LastActivity: session.LastActivity,
		ExpiresAt:    session.ExpiresAt,
		IpAddress:    sql.NullString{String: session.IpAddress, Valid: session.IpAddress != ""},
		UserAgent:    sql.NullString{String: session.UserAgent, Valid: session.UserAgent != ""},
		IsActive:     session.IsActive,
		RevokedAt:    revokedAt,
		OtpVerified:  sql.NullBool{Bool: session.OTPVerified, Valid: true},
		OtpExpiresAt: sql.NullTime{Time: session.OtpExpiresAt, Valid: true},
		Otp:          sql.NullString{String: session.Otp, Valid: session.Otp != ""},
		OtpAttempts:  sql.NullInt32{Int32: int32(session.OtpAttempts), Valid: true},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateSession(ctx context.Context, session *entity.Session) error {
	// Convert sessionID from string to UUID
	UserUUID, err := uuid.Parse(session.UserID.String())
	if err != nil {
		return fmt.Errorf("invalid session ID format: %w", err)
	}

	// Call UpdateSession with the mapped parameters
	_, err = r.store.UpdateSession(ctx, db.UpdateSessionParams{
		UserID:       UserUUID,
		LastActivity: session.LastActivity,
		ExpiresAt:    session.ExpiresAt,
		IpAddress:    sql.NullString{String: session.IpAddress, Valid: session.IpAddress != ""},
		UserAgent:    sql.NullString{String: session.UserAgent, Valid: session.UserAgent != ""},
		DeviceInfo:   pqtype.NullRawMessage{RawMessage: session.DeviceInfo.RawMessage, Valid: session.DeviceInfo != nil},
		IsActive:     session.IsActive,
		RevokedAt:    sql.NullTime{Time: *session.RevokedAt, Valid: true},
	},
	)

	return err
}

// DeleteSession deletes a user session by its ID.
func (r *UserRepository) DeleteSession(ctx context.Context, id string) error {
	// Convert sessionID from string to UUID
	sessionUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid session ID format: %w", err)
	}

	// Call DeleteSession with the mapped parameters
	err = r.store.DeleteSession(ctx, sessionUUID)

	return err
}

func (r *UserRepository) UpdateOtp(ctx context.Context, userOtp *entity.UpdateOtp) error {
	// Get User ID from the session
	user, err := r.GetUserByEmail(ctx, userOtp.Email)
	if err != nil {
		return err
	}

	// Call UpdateOtp with the mapped parameters
	_, err = r.store.UpdateSession(ctx, db.UpdateSessionParams{
		UserID:       user.ID,
		Otp:          sql.NullString{String: userOtp.Otp, Valid: true},
		OtpExpiresAt: sql.NullTime{Time: userOtp.OtpExpiresAt, Valid: true},
		OtpAttempts:  sql.NullInt32{Int32: int32(userOtp.OtpAttempts), Valid: true},
		OtpVerified:  sql.NullBool{Bool: userOtp.OTPVerified, Valid: true},
	})
	if err != nil {
		return err
	}

	return nil

}

func (r *UserRepository) GetOtp(ctx context.Context, email string) (*entity.UpdateOtp, error) {
	// Retrieve user by email
	user, err := r.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user by email %s: %w", email, err)
	}

	// Convert userID from string to UUID
	userUUID, err := uuid.Parse(string(user.ID.String()))
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Call GetOtp with the mapped parameters
	otp, err := r.store.GetSession(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	// Map session details to the entity.Session struct
	return &entity.UpdateOtp{
		Otp:          otp.Otp.String,
		OtpExpiresAt: otp.OtpExpiresAt.Time,
		OtpAttempts:  int(otp.OtpAttempts.Int32),
	}, nil

}
