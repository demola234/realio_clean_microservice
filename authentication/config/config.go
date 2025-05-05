package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	Environment       string `mapstructure:"ENVIRONMENT"`

	// Google OAuth
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`

	// Facebook OAuth
	FacebookAppID     string `mapstructure:"FACEBOOK_APP_ID"`
	FacebookAppSecret string `mapstructure:"FACEBOOK_APP_SECRET"`

	// Apple OAuth
	AppleClientID     string `mapstructure:"APPLE_CLIENT_ID"`
	AppleClientSecret string `mapstructure:"APPLE_CLIENT_SECRET"`
	AppleTeamID       string `mapstructure:"APPLE_TEAM_ID"`
	AppleKeyID        string `mapstructure:"APPLE_KEY_ID"`
	ApplePrivateKey   string `mapstructure:"APPLE_PRIVATE_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName(".env")

	// Set the path to look for the configuration file
	viper.AddConfigPath(path)

	// Set the file type to .env
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_SOURCE", "postgres://root:secret@localhost:5433/defi?sslmode=disable")
	viper.SetDefault("GRPC_SERVER_ADDRESS", ":50051")
	viper.SetDefault("HTTP_SERVER_ADDRESS", ":8080")
	viper.SetDefault("TOKEN_SYMMETRIC_KEY", "12345678901234567890123456789012")
	viper.SetDefault("Environment", "development")

	// GoogleOAuth
	viper.SetDefault("GOOGLE_CLIENT_ID", "123456789012-abcdefghijklmnopqrstuvwxyz.apps.googleusercontent.com")
	viper.SetDefault("GOOGLE_CLIENT_SECRET", "abcdefghijklmnopqrstuvwxyz123456")

	// Facebook OAuth
	viper.SetDefault("FACEBOOK_APP_ID", "123456789012345")
	viper.SetDefault("FACEBOOK_APP_SECRET", "abcdefghijklmnopqrstuvwxyz123456")

	// Apple OAuth
	viper.SetDefault("APPLE_CLIENT_ID", "com.example.app")
	viper.SetDefault("APPLE_CLIENT_SECRET", "abcdefghijklmnopqrstuvwxyz123456")
	viper.SetDefault("APPLE_TEAM_ID", "ABCDEFGHIJK")
	viper.SetDefault("APPLE_KEY_ID", "ABCDEFGHIJK")
	viper.SetDefault("APPLE_PRIVATE_KEY", "-----BEGIN PRIVATE KEY-----\nabcdefghijklmnopqrstuvwxyz123456\n-----END PRIVATE KEY-----")

	viper.AutomaticEnv()

	// Set the type of the configuration file
	viper.SetConfigType("env")

	// Read the configuration file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Unmarshal the configuration into the Config struct
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
