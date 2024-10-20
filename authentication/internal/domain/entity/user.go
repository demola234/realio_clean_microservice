package entity

import (
	"time"
)

type User struct {
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	AccessToken       string    `json:"access_token"`
	RefreshToken      string    `json:"refresh_token"`
}
