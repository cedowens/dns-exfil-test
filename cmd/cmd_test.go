package cmd_test

import (
	"os"
	"testing"

	"github.com/dns-tun/cmd"

	"github.com/reedobrien/checkers"
)

const (
	defaultListenPort = "5353"
	defaultListenAddr = "0.0.0.0"

	listenAddrEnvVar = "LISTEN_ADDR_ENV_VAR"
	listenPortEnvVar = "LISTEN_PORT_ENV_VAR"
)

// TestGetDefaultValue tests we receive the default value when the env var isn't set.
func TestGetDefaultValue(t *testing.T) {
	table := []struct {
		name string
		got  string
		want string
	}{
		{"defaultListPort", cmd.GetEnvValue(listenPortEnvVar, defaultListenPort), defaultListenPort},
		{"defaultListAddr", cmd.GetEnvValue(listenAddrEnvVar, defaultListenAddr), defaultListenAddr},
	}

	for _, tc := range table {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			checkers.Equals(t, tc.got, tc.want)
		})
	}
}

// TestDetEnvVarValue tests that we receive the env var defined value.
func TestDetEnvVarValue(t *testing.T) {
	err := os.Setenv(listenPortEnvVar, "12345")
	checkers.OK(t, err)
	err = os.Setenv(listenAddrEnvVar, "1.1.1.1")
	checkers.OK(t, err)

	table := []struct {
		name string
		got  string
		want string
	}{
		{"envPortValue", cmd.GetEnvValue(listenPortEnvVar, defaultListenPort), "12345"},
		{"envListAddr", cmd.GetEnvValue(listenAddrEnvVar, defaultListenAddr), "1.1.1.1"},
	}

	for _, tc := range table {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			checkers.Equals(t, tc.got, tc.want)
		})
	}

	err = os.Unsetenv(listenPortEnvVar)
	checkers.OK(t, err)

	err = os.Unsetenv(listenAddrEnvVar)
	checkers.OK(t, err)
}
