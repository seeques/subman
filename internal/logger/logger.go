package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(level string) zerolog.Logger {
    // Parse level
    lvl, err := zerolog.ParseLevel(level)
    if err != nil {
        lvl = zerolog.InfoLevel
    }
    zerolog.SetGlobalLevel(lvl)

    return zerolog.New(os.Stdout).With().Timestamp().Logger()
}