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
	return zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)
}