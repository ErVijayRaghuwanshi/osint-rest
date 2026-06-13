package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger

func New() Logger {
	level := zerolog.DebugLevel
	if l := os.Getenv("LOG_LEVEL"); l != "" {
		parsed, err := zerolog.ParseLevel(l)
		if err == nil {
			level = parsed
		}
	}

	if os.Getenv("GIN_MODE") != "release" {
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}).
			With().
			Timestamp().
			Logger().
			Level(level)
	}

	return zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)
}