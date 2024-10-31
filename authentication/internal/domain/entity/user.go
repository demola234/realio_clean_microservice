package entity

import (
	"time"

	"github.com/google/uuid"
)

// User entity based on the users table schema
type User struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
