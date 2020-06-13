package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/dns-tun/cmd"
)

const (
	appName           = "dnsTun"
	defaultListenPort = "53"
	defaultListenAddr = "0.0.0.0"

	listenAddrEnvVar = "LISTEN_ADDR_ENV_VAR"
	listenPortEnvVar = "LISTEN_PORT_ENV_VAR"
)

type dnsServer struct {
	ctx    context.Context
	wg     *sync.WaitGroup
	logger zerolog.Logger
	lstnr  *net.UDPConn
}

func (srvr *dnsServer) requestHandler() {

}
func (srvr *dnsServer) runServer(stop chan os.Signal) {

	buf := make([]byte, 1024)
	_, addr, err := srvr.lstnr.ReadFromUDP(buf)
	if err != nil {
		srvr.logger.Error().Err(err).Msg("error reading from UDPConff")
	}
	fmt.Printf("%+v\t\t%+v", string(buf), addr)

	// block until we receive a stop signal
	<-stop

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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	srvr := &dnsServer{
		ctx:    context.Background(),
		wg:     &sync.WaitGroup{},
		logger: cmd.SetupLogger(*debug, appName),
	}

	iport, err := strconv.Atoi(*port)
	if err != nil {
		srvr.logger.Fatal().Err(err).Msg("unable to convert port to int")
	}

	udpAddrPort := &net.UDPAddr{
		IP:   net.ParseIP(*address),
		Port: iport,
	}

	lstnr, err := net.ListenUDP("udp", udpAddrPort)
	if err != nil {
		srvr.logger.Fatal().Err(err).Msg("unable to create listener")
	}

	srvr.lstnr = lstnr

	srvr.runServer(stop)

}
