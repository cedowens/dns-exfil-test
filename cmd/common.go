package cmd

import (
	"os"

	"github.com/reedobrien/logging"
	"github.com/rs/zerolog"
)

func GetEnvValue(envvar, defValue string) string {
	val, ok := os.LookupEnv(envvar)
	if ok {
		return val
	}
	return defValue
}

func SetupLogger(debug bool, appName string) zerolog.Logger {
	logger := logging.NewLogger(appName, debug, nil)

	return logger
}
