package logging

import (
	"io"
	stdliblog "log"
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NewLogger returns a constructed logger.
func NewLogger(app string, verbose bool, out io.Writer) zerolog.Logger {
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	// Show nanoseconds.
	zerolog.TimeFieldFormat = TimeFormat

	// NanoSeconds? Default is milliseconds.
	// zerolog.DurationFieldUnit = time.Nanosecond
	// zerolog.DurationFieldInteger = true

	// Setup a logger.
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get hostname")
	}

	if out == nil {
		out = os.Stdout
	}

	logger := zerolog.New(out).With().
		Timestamp().
		Str("app", app).
		Str("app_host", hostname).
		Str("go_version", runtime.Version()).
		Logger()

	// Tell stdlib tools that use the default logger to use our logger.
	stdliblog.SetFlags(0)
	stdliblog.SetOutput(logger)

	return logger
}
