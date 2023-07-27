package handler

import (
	"net"
	"os"
	"strings"

	"github.com/miekg/dns"
	"github.com/rs/zerolog"
)

type DNSHandler struct {
	Logger     zerolog.Logger
	Outputfile string
}

func New(logger zerolog.Logger) *DNSHandler {
	return &DNSHandler{
		Logger:     logger,
		Outputfile: "outfile",
	}
}

func (dh *DNSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name

		domainsToAddresses := map[string]string{
			"domain": "0.0.0.0",
		}

		address, ok := domainsToAddresses[domain]

		if ok {
			dh.Logger.Info().Msg("[+] Receiving bytes...")

			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})

			split := strings.Split(domain, ".")
			split2 := split[0]

			p, err := os.OpenFile(dh.Outputfile,
				os.O_WRONLY|os.O_APPEND|os.O_CREATE,
				0o666)
			if err != nil {
				dh.Logger.Err(err).Msg("[-] Error creating the outfile...")

				return
			}

			defer p.Close()

			_, err = p.Write([]byte(split2))
			if err != nil {
				dh.Logger.Err(err).Msg("error writing to outfile")

				return
			}
		}
	}
	err := w.WriteMsg(&msg)
	if err != nil {
		dh.Logger.Err(err).Msg("error calling WriteMsg")
	}
}
