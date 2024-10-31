package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port                      string `mapstructure:"API_GATEWAY_PORT"`
	AuthenticationGRPCAddress string `mapstructure:"AUTH_GRPC_ADDRESS"`
	SentryConfig              string `mapstructure:"SENTRY_CONFIG"`
	SentryConst               string `mapstructure:"SENTRY_CONST"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
