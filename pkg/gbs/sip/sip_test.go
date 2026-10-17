package sip

import (
	"testing"
)

func TestParseHeader(t *testing.T) {
	const s = `Contact: sip:1678@80.79.5.134;expires=3600`
	h, err := ParseHeader(s)
	if err != nil {
		t.Fatalf("ParseHeader failed: %v", err)
	}
	t.Log(h)
}
