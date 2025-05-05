package repository

import (
	"context"

	domain "github.com/demola234/authentication/internal/domain/entity"
)

// OAuthRepository defines the repository contract for OAuth-related operations.
type OAuthRepository interface {
	// ValidateGoogleToken validates a Google OAuth token.
	ValidateGoogleToken(ctx context.Context, token string) (*domain.OAuthUserInfo, error)
	// ValidateFacebookToken validates a Facebook OAuth token.
	ValidateFacebookToken(ctx context.Context, token string) (*domain.OAuthUserInfo, error)
	// ValidateAppleToken validates an Apple OAuth token.
	ValidateAppleToken(ctx context.Context, token string) (*domain.OAuthUserInfo, error)
	// ValidateProviderToken validates an OAuth token for a specific provider.
	ValidateProviderToken(ctx context.Context, provider, token string) (*domain.OAuthUserInfo, error)
}