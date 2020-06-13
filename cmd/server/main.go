package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/dns-tun/cmd"
)

const (
	appName           = "dnsTun"
	defaultListenPort = "53"
	defaultListenAddr = "0.0.0.0"

	clientSrcAddrLogField = "src_addr"
	clientBytesReceived   = "bytes_received"

	udpDataBufferSize = 1024

	listenAddrEnvVar = "LISTEN_ADDR_ENV_VAR"
	listenPortEnvVar = "LISTEN_PORT_ENV_VAR"
)

type dnsServer struct {
	ctx     context.Context
	wg      *sync.WaitGroup
	logger  zerolog.Logger
	udpaddr *net.UDPAddr
}

func (srvr *dnsServer) msgHandler(udpconnn *net.UDPConn) {
	for {
		// This is most definitely not performant...should create once and rezero.
		udpbuf := make([]byte, udpDataBufferSize)
		n, addr, err := udpconnn.ReadFromUDP(udpbuf)

		if err != nil {
			srvr.logger.Error().Err(err).Msg("error reading the udp connection")
		}

		srvr.logger.Debug().
			Str(clientSrcAddrLogField, addr.String()).
			Int(clientBytesReceived, n).
			Msg("processing message")

		fmt.Printf("\t\tString: %+v\n", string(udpbuf))
	}
}
func (srvr *dnsServer) runServer(stop chan os.Signal) {
	udpconn, err := net.ListenUDP("udp", srvr.udpaddr)
	if err != nil {
		srvr.logger.Fatal().Err(err).Msg("error when getting udpconn from net.ListenUDP")
	}

	go func() {
		for {
			srvr.msgHandler(udpconn)
		}
	}()

	<-stop
	udpconn.Close() //nolint:errcheck,gosec
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

	udpAddrPort, err := net.ResolveUDPAddr("udp", *address+":"+*port)
	if err != nil {
		srvr.logger.Fatal().Err(err).Msg("unable to create UDPAddress from address and port")
	}

	srvr.udpaddr = udpAddrPort

	srvr.runServer(stop)
	srvr.logger.Info().Msg("exiting")
}
