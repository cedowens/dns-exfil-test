package main

import (
	"flag"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/dns-tun/cmd"
)

const (
	appName           = "dnsTun"
	defaultListenPort = "5353"
	defaultListenAddr = "0.0.0.0"

	listenAddrEnvVar = "LISTEN_ADDR_ENV_VAR"
	listenPortEnvVar = "LISTEN_PORT_ENV_VAR"
)

type dnsServer struct {
	logger zerolog.Logger
}

func main() {
	var (
		address = flag.String("lstnAddr",
			cmd.GetEnvValue(listenAddrEnvVar, defaultListenAddr),
			"What address should the server listen on?")
		port = flag.String("lstnPort",
			cmd.GetEnvValue(listenPortEnvVar, defaultListenPort),
			"What port should we listen on?")
		debug = flag.Bool("debug",
			false,
			"Run the server in debug")
	)

	flag.Parse()

	srvr := &dnsServer{
		logger: cmd.SetupLogger(*debug, appName),
	}

	fmt.Printf("%+v\t%+v\t%+v\n", srvr, address, port)

}
