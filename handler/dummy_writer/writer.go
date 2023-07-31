package dummywriter

import (
	"bytes"
	"net"

	"github.com/miekg/dns"
)

type DW struct {
	GotMsg   *dns.Msg
	GotBytes *bytes.Buffer
}

func New(opts ...func(*DW)) *DW {
	dw := &DW{}

	for _, opt := range opts {
		opt(dw)
	}

	return dw
}

func (dw *DW) LocalAddr() net.Addr {
	return &net.IPAddr{}
}

func (dw *DW) RemoteAddr() net.Addr {
	return &net.IPAddr{}
}

func (dw *DW) WriteMsg(m *dns.Msg) error {
	dw.GotMsg = m

	return nil
}

func (dw *DW) Write(b []byte) (int, error) {
	return dw.GotBytes.Write(b)
}

func (dw *DW) Close() error {
	return nil
}

func (dw *DW) TsigStatus() error {
	return nil
}

func (dw *DW) TsigTimersOnly(b bool) {
	return
}

func (dw *DW) Hijack() {
	return
}
