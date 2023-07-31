package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cedowens/dns-exfil-test/handler"
	"github.com/miekg/dns"
	"github.com/rs/zerolog"
)

const drainTime = 10 * time.Second

func main() {
	listenaddr := flag.String("listenaddr", ":53", "address and port to listen on")
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	srv := &dns.Server{Addr: *listenaddr, Net: "udp"}

	logger := zerolog.New(os.Stdout)

	srv.Handler = handler.New(logger)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), drainTime)
	defer cancel()

	_ = srv.ShutdownContext(ctx)
}
