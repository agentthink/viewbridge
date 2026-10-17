package onvifserver

import (
	"testing"

	"github.com/gowvp/owl/internal/conf"
)

func TestAdvertiseHostPrefersSDPIP(t *testing.T) {
	cfg := &conf.Bootstrap{
		Media: conf.Media{IP: "127.0.0.1", SDPIP: "192.168.1.244"},
		Sip:   conf.SIP{Host: "192.168.1.241"},
	}
	if got := advertiseHost(cfg); got != "192.168.1.244" {
		t.Fatalf("got %q want 192.168.1.244", got)
	}
}
