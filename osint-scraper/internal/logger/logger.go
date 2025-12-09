package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger

func New() Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
