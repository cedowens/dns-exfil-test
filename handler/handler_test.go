package handler_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/cedowens/dns-exfil-test/handler"
	dummywriter "github.com/cedowens/dns-exfil-test/handler/dummy_writer"
	"github.com/miekg/dns"
	"github.com/rs/zerolog"

	"github.com/dvterry/gotest"
)

func TestServeDNSNoDomainMatch(t *testing.T) {
	gotMsg := &dns.Msg{}

	dw := dummywriter.New(func(dw *dummywriter.DW) { dw.GotMsg = gotMsg })

	logBucket := &bytes.Buffer{}
	logger := zerolog.New(logBucket)

	tut := handler.New(logger)
	testMsg := &dns.Msg{
		Question: []dns.Question{
			{
				Qtype: dns.TypeA,
				Name:  "thisisnotinthemap",
			},
		},
	}

	tut.ServeDNS(dw, testMsg)

	// We expect an empty *dns.Msg
	gotest.IsEqual(t, gotMsg, &dns.Msg{})
	// We shouldn't have an error log in the buffer
	gotest.IsEqual(t, logBucket.String(), "")
}

func TestServeDNSDomainMatch(t *testing.T) {
	gotBytes := &bytes.Buffer{}
	dw := dummywriter.New(func(dw *dummywriter.DW) { dw.GotBytes = gotBytes })

	logBucket := &bytes.Buffer{}
	logger := zerolog.New(logBucket)

	tut := handler.New(logger)
	testMsg := &dns.Msg{
		Question: []dns.Question{
			{
				Qtype: dns.TypeA,
				Name:  "domain",
			},
		},
	}

	tut.ServeDNS(dw, testMsg)

	fmt.Printf("%s\n", gotBytes.String())

	fmt.Printf("%s\n", logBucket.String())
}
