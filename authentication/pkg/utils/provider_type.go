package utils

import (
	"strings"
)

// ProviderType represents authentication provider information
type ProviderType struct {
	Name      string `json:"provider"`      // The provider name: "local", "google", "facebook", etc.
	ID        string `json:"provider_id"`   // Provider-specific user ID
	TokenData string `json:"token_data"`    // Store additional token data (optional)
	
	// For backward compatibility with existing code
	Email    string `json:"email,omitempty"`    // Local email (will be moved to ID)
	Google   string `json:"google,omitempty"`   // Google ID (deprecated)
	Facebook string `json:"facebook,omitempty"` // Facebook ID (deprecated)
	Twitter  string `json:"twitter,omitempty"`  // Twitter ID (deprecated)
	Github   string `json:"github,omitempty"`   // Github ID (deprecated)
	Local    string `json:"local,omitempty"`    // Local ID (deprecated)
}

// Helper methods
func (p *ProviderType) IsLocal() bool {
	return p.Name == "local" || p.Name == ""
}

func (p *ProviderType) IsOAuth() bool {
	return !p.IsLocal()
}

// GetProviderInfo returns the provider name and ID
func (p *ProviderType) GetProviderInfo() (provider string, providerID string) {
	// First check the new format
	if p.Name != "" {
		return p.Name, p.ID
	}
	
	// Backward compatibility: check old format
	if p.Google != "" {
		return "google", p.Google
	}
	if p.Facebook != "" {
		return "facebook", p.Facebook
	}
	if p.Twitter != "" {
		return "twitter", p.Twitter
	}
	if p.Github != "" {
		return "github", p.Github
	}
	if p.Email != "" {
		return "local", p.Email
	}
	if p.Local != "" {
		return "local", p.Local
	}
	
	// Default to local
	return "local", ""
}

// SetProvider sets the provider information using the new format
func (p *ProviderType) SetProvider(provider, id string) {
	p.Name = strings.ToLower(provider)
	p.ID = id
	
	// Clear old format fields for consistency
	p.Email = ""
	p.Google = ""
	p.Facebook = ""
	p.Twitter = ""
	p.Github = ""
	p.Local = ""
}

// Backward compatibility getters
func (p *ProviderType) GetEmail() string {
	if p.Name == "local" {
		return p.ID
	}
	return p.Email
}

func (p *ProviderType) GetGoogle() string {
	if p.Name == "google" {
		return p.ID
	}
	return p.Google
}

func (p *ProviderType) GetFacebook() string {
	if p.Name == "facebook" {
		return p.ID
	}
	return p.Facebook
}

func (p *ProviderType) GetTwitter() string {
	if p.Name == "twitter" {
		return p.ID
	}
	return p.Twitter
}

func (p *ProviderType) GetGithub() string {
	if p.Name == "github" {
		return p.ID
	}
	return p.Github
}

func (p *ProviderType) GetLocal() string {
	if p.IsLocal() {
		return p.ID
	}
	return p.Local
}