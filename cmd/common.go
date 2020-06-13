// Package cmd implements helper functions for the entry points.
package cmd

import (
	"os"

	"github.com/reedobrien/logging"
	"github.com/rs/zerolog"
)

// GetEnvValue accepts a default value and env var and returns value within
// env var if present, otherwise returns default.
func GetEnvValue(envvar, defValue string) string {
	val, ok := os.LookupEnv(envvar)
	if ok {
		return val
	}

	return defValue
}

// SetupLogger returns an instantiated zerolog.Logger.
func SetupLogger(debug bool, appName string) zerolog.Logger {
	logger := logging.NewLogger(appName, debug, nil)

	return logger
}
