package entity

import (
	"errors"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
	ErrProviderInvalid = errors.New("invalid oauth provider")
)

// OAuthUserInfo represents the user information from OAuth providers
type OAuthUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
	Provider      string `json:"provider"`
}
