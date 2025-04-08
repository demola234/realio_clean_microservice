package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName(".env")

	// Set the path to look for the configuration file
	viper.AddConfigPath(path)

	// Set the file type to .env
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_SOURCE", "postgres://root:secret@localhost:5433/defi?sslmode=disable")
	viper.SetDefault("GRPC_SERVER_ADDRESS", ":50051")
	viper.SetDefault("TOKEN_SYMMETRIC_KEY", "12345678901234567890123456789012")
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
