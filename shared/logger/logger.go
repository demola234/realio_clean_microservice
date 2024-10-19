package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	
)

// InitLogger initializes the global logger for the application
func InitLogger() {
	// Set the global log level (debug, info, warn, error)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Configure zerolog with human-readable time format for better readability
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	// Optional: Use human-friendly logging in development environments
	if os.Getenv("APP_ENV") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	log.Info().Msg("Logger initialized successfully")
}

// SetLogLevel dynamically changes the global log level
func SetLogLevel(level string) {
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Str("level", level).Msg("Log level set")
}
