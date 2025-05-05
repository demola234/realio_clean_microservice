package entity

import (
	"time"

	"github.com/demola234/authentication/pkg/utils"
	"github.com/google/uuid"
)

// User entity based on the users table schema
type User struct {
	ID             uuid.UUID          `json:"id"`
	FullName       string             `json:"name"`
	Email          string             `json:"email"`
	Bio            string             `json:"bio"`
	Provider       utils.ProviderType `json:"provider"`
	Username       string             `json:"username"`
	ProviderID     string             `json:"provider_id"`
	ProfilePicture string             `json:"profile_picture"`
	Password       string             `json:"password"`
	Role           string             `json:"role"`
	Phone          string             `json:"phone"`
	EmailVerified  bool               `json:"email_verified"`
	IsActive       bool               `json:"is_active"`
	LastLogin      time.Time          `json:"last_login"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}
