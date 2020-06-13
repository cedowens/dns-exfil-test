package cmd_test

import (
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

// TestGetDefaultValue tests we receive the default value when the env var isn't set
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
