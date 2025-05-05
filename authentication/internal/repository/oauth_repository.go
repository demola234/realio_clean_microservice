package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/demola234/authentication/config"
	domain "github.com/demola234/authentication/internal/domain/entity"
	repo "github.com/demola234/authentication/internal/domain/repository"
)

// OAuthRepositoryImpl implements the OAuthRepository interface
type OAuthRepositoryImpl struct {
	config *config.Config
	client *http.Client
}

// NewOAuthRepository creates a new OAuth repository
func NewOAuthRepository(cfg *config.Config) repo.OAuthRepository {
	return &OAuthRepositoryImpl{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *OAuthRepositoryImpl) ValidateProviderToken(ctx context.Context, provider, token string) (*domain.OAuthUserInfo, error) {
	switch strings.ToLower(provider) {
	case "google":
		return r.ValidateGoogleToken(ctx, token)
	case "facebook":
		return r.ValidateFacebookToken(ctx, token)
	case "apple":
		return r.ValidateAppleToken(ctx, token)
	default:
		return nil, domain.ErrProviderInvalid
	}
}

// ValidateGoogleToken validates a Google OAuth token.
func (r *OAuthRepositoryImpl) ValidateGoogleToken(ctx context.Context, token string) (*domain.OAuthUserInfo, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", token)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to validate token: %s", resp.Status)
	}

	var googleResult struct {
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Sub           string `json:"sub"`
		EmailVerified string `json:"email_verified"`
		Exp           string `json:"exp"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if googleResult.Exp != "" {
		expTime, err := time.Parse(time.RFC3339, googleResult.Exp)
		if err == nil && expTime.Before(time.Now()) {
			return nil, domain.ErrTokenExpired
		}
	}

	return &domain.OAuthUserInfo{
		Email:         googleResult.Email,
		Name:          googleResult.Name,
		Picture:       googleResult.Picture,
		ID:            googleResult.Sub,
		EmailVerified: googleResult.EmailVerified == "true",
		Provider:      "google",
	}, nil
}

// ValidateFacebookToken validates a Facebook OAuth token.
func (r *OAuthRepositoryImpl) ValidateFacebookToken(ctx context.Context, token string) (*domain.OAuthUserInfo, error) {
	url := fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%s&access_token=%s|%s",
		token, r.config.FacebookAppID, r.config.FacebookAppSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to validate token: %s", resp.Status)
	}

	var facebookResult struct {
		Data struct {
			AppID     string `json:"app_id"`
			ExpiresAt int64  `json:"expires_at"`
			IsValid   bool   `json:"is_valid"`
			Scopes    []struct {
				Scope     string `json:"scope"`
				Valid     bool   `json:"valid"`
				ErrorName string `json:"error_name"`
			} `json:"scopes"`
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&facebookResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !facebookResult.Data.IsValid {
		return nil, domain.ErrInvalidToken
	}

	if !facebookResult.Data.IsValid {
		return nil, domain.ErrInvalidToken
	}

	if facebookResult.Data.AppID != r.config.FacebookAppID {
		return nil, domain.ErrInvalidToken
	}

	// Check expiration
	if facebookResult.Data.ExpiresAt != 0 && time.Now().Unix() > facebookResult.Data.ExpiresAt {
		return nil, domain.ErrTokenExpired
	}

	userInfo, err := r.getFaceBookUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}

	userInfo.ID = facebookResult.Data.UserID
	userInfo.Provider = "facebook"

	return userInfo, nil
}

func (r *OAuthRepositoryImpl) getFaceBookUserInfo(ctx context.Context, token string) (*domain.OAuthUserInfo, error) {
	url := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email,picture&access_token=%s", token)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &domain.OAuthUserInfo{
		ID:      userInfo.ID,
		Name:    userInfo.Name,
		Email:   userInfo.Email,
		Picture: userInfo.Picture.Data.URL,
	}, nil
}

func (r *OAuthRepositoryImpl) ValidateAppleToken(ctx context.Context, token string) (*domain.OAuthUserInfo, error) {
	// Apple's validation endpoint
	data := url.Values{}
	data.Set("client_id", r.config.AppleClientID)
	data.Set("client_secret", r.config.AppleClientSecret)
	data.Set("token", token)
	data.Set("token_type_hint", "id_token")

	url := "https://appleid.apple.com/auth/token"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to validate token: %s", resp.Status)
	}

	var appleResult struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&appleResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse the ID token to get user info
	userInfo, err := r.parseAppleIDToken(appleResult.IDToken)
	if err != nil {
		return nil, err
	}

	userInfo.Provider = "apple"

	return userInfo, nil
}

func (r *OAuthRepositoryImpl) parseAppleIDToken(idToken string) (*domain.OAuthUserInfo, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Exp           int64  `json:"exp"`
	}

	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	if time.Now().Unix() > claims.Exp {
		return nil, domain.ErrTokenExpired
	}

	var result = &domain.OAuthUserInfo{
		ID:            claims.Sub,
		Email:         claims.Email,
		Name:          "",
		Picture:       "",
		EmailVerified: claims.EmailVerified == "true",
	}

	return result, nil

}
