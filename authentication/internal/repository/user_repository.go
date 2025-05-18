package repository

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	token "github.com/demola234/api_gateway/infrastructure/middleware/token_maker"
	"github.com/demola234/authentication/config"
	db "github.com/demola234/authentication/db/sqlc" // SQLC generated code for interacting with the database
	"github.com/demola234/authentication/internal/domain/entity"
	"github.com/demola234/authentication/pkg/utils"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// UserRepository implements the AuthRepository interface.
// This struct interacts with the database using SQLC-generated code.
type UserRepository struct {
	store db.Store
}

// UpdateUser implements repository.UserRepository.
func (r *UserRepository) UploadProfileImage(ctx context.Context, content io.Reader, userId uuid.UUID) (string, error) {
	// Create a new Cloudinary client
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	// Upload the image to the cloudinary
	cld, _ := cloudinary.NewFromURL(configs.CloudinaryUrl)

	// Get the preferred name of the file if its not supplied
	result, err := cld.Upload.Upload(ctx, content, uploader.UploadParams{
		PublicID: userId.String(),
		Tags:     strings.Split(",", userId.String()),
	})

	if err != nil {
		return "", err
	}

	updateUser := db.UpdateUserProfilePictureParams{
		ID: userId,
		ProfilePicture: sql.NullString{
			String: result.SecureURL,
			Valid:  true,
		},
	}

	return updateUser.ProfilePicture.String, nil

}

// UpdateUser implements repository.UserRepository.
func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {

	// Check if the user exists
	updateUser := db.UpdateUserParams{
		ID:       user.ID,
		Name:     user.FullName,
		Username: user.Username,
		ProfilePicture: sql.NullString{
			String: user.ProfilePicture,
			Valid:  true,
		},
		Bio: sql.NullString{
			String: user.Bio,
			Valid:  true,
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  true,
		},
		Role: sql.NullString{
			String: user.Role,
			Valid:  true,
		},
		Email: user.Email,
		Password: sql.NullString{
			String: user.Password,
			Valid:  true,
		},
	}

	updatedUser, err := r.store.UpdateUser(ctx, updateUser)
	if err != nil {
		return err
	}

	// Update the user in the cache
	user.ID = updatedUser.ID
	user.FullName = updatedUser.Name
	user.Username = updatedUser.Username
	user.ProfilePicture = updatedUser.ProfilePicture.String
	user.IsActive = updatedUser.IsActive.Bool
	user.EmailVerified = updatedUser.EmailVerified.Bool
	user.Phone = updatedUser.Phone.String
	user.Provider = utils.ProviderType{
		Name:      updatedUser.Provider.String,
		ID:        updatedUser.ProviderID.String,
		TokenData: updatedUser.ProviderID.String,
	}
	user.EmailVerified = updatedUser.EmailVerified.Bool

	return nil

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

	password := userDetails.Password.String

	// Create ProviderType using the new unified format
	provider := utils.ProviderType{}
	if userDetails.Provider.Valid {
		provider.SetProvider(userDetails.Provider.String, userDetails.ProviderID.String)
	} else {
		// Default to local if provider not set
		provider.SetProvider("local", userDetails.Email)
	}

	return &entity.User{
		ID:             userDetails.ID,
		FullName:       userDetails.Name,
		Username:       userDetails.Username,
		Email:          userDetails.Email,
		Provider:       provider,
		ProviderID:     userDetails.ProviderID.String,
		ProfilePicture: userDetails.ProfilePicture.String,
		Role:           userDetails.Role.String,
		Password:       password,
		Phone:          userDetails.Phone.String,
		EmailVerified:  userDetails.EmailVerified.Bool,
		IsActive:       userDetails.IsActive.Bool,
		LastLogin:      userDetails.LastLogin.Time,
		CreatedAt:      userDetails.CreatedAt.Time,
		UpdatedAt:      userDetails.UpdatedAt.Time,
	}, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	// Get provider information using the unified method
	provider, providerID := user.Provider.GetProviderInfo()

	// Convert to sql.NullString
	var providerIDNullable sql.NullString
	if providerID != "" {
		providerIDNullable = sql.NullString{String: providerID, Valid: true}
	} else {
		providerIDNullable = sql.NullString{Valid: false}
	}

	// Hash password only for local users
	var hashedPasswordString sql.NullString
	if provider == "local" && user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		hashedPasswordString = sql.NullString{String: hashedPassword, Valid: true}
	} else {
		hashedPasswordString = sql.NullString{Valid: false}
	}

	username := user.Username
	if username == "" {
		if strings.Contains(user.Email, "@") {
			username = user.Email[:strings.Index(user.Email, "@")]
		} else {
			username = user.Email // Fallback if email format is invalid
		}
	}

	_, err := r.store.CreateUser(ctx, db.CreateUserParams{
		ID:             user.ID,
		Name:           user.FullName,
		Username:       username,
		Email:          user.Email,
		Password:       hashedPasswordString,
		ProfilePicture: sql.NullString{String: user.ProfilePicture, Valid: user.ProfilePicture != ""},
		Bio:            sql.NullString{String: "", Valid: true},
		Role:           sql.NullString{String: user.Role, Valid: user.Role != ""},
		Phone:          sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		Provider:       sql.NullString{String: provider, Valid: true},
		ProviderID:     providerIDNullable,
		EmailVerified:  sql.NullBool{Bool: false, Valid: true},
		IsActive:       sql.NullBool{Bool: true, Valid: true},
		LastLogin:      sql.NullTime{Valid: false},
	})

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// UpdatePassword updates a user's password in the database.
func (r *UserRepository) UpdatePassword(ctx context.Context, userID string, newPassword string) error {

	// Hash the new password
	newPass := sql.NullString{String: newPassword, Valid: true}

	_, err := r.store.UpdateUser(ctx, db.UpdateUserParams{
		Password: newPass,
	})

	return err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	userDetails, err := r.store.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	password := userDetails.Password.String
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty for user with ID %s", id)
	}

	return &entity.User{
		ID:        userDetails.ID,
		FullName:  userDetails.Name,
		Email:     userDetails.Email,
		CreatedAt: userDetails.CreatedAt.Time,
		Password:  password,
		Role:      userDetails.Role.String,
		Phone:     userDetails.Phone.String,
		UpdatedAt: userDetails.UpdatedAt.Time,
	}, nil

}

func (r *UserRepository) GetUserSession(ctx context.Context, sessionID uuid.UUID) (*entity.Session, error) {

	// Retrieve session details from the store
	sessionDetails, err := r.store.GetSessionByUserID(ctx, sessionID)
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
	otp, err := r.store.GetSessionByUserID(ctx, userUUID)
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
