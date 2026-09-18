// Package logger provides a wrapper around the zerolog logging library for structured logging.
package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
)

// New initializes and returns a new zerolog.Logger instance configured for the application.
func New() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	if os.Getenv("GIN_MODE") == "release" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}
	return log.Logger

}
