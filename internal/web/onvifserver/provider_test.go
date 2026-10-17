package onvifserver

import (
	"strings"
	"testing"

	"github.com/gowvp/owl/internal/core/ipc"
)

func TestProfileName(t *testing.T) {
	p := &StreamProvider{
		devices: map[string]deviceMeta{"dev1": {name: "NVR一楼", online: true}},
	}
	if got := p.profileName(&ipc.Channel{ID: "ch1", DID: "dev1", Name: "前门", IsOnline: true}); got != "前门-NVR一楼" {
		t.Fatalf("got %q", got)
	}
	p.devices["dev2"] = deviceMeta{online: true}
	if got := p.profileName(&ipc.Channel{ID: "ch2", DID: "dev2", ChannelID: "34010000001320000001", IsOnline: true}); got != "34010000001320000001-dev2" {
		t.Fatalf("got %q", got)
	}
	p.devices["gb-dev-1"] = deviceMeta{name: "门口摄像机", online: true}
	if got := p.profileName(&ipc.Channel{
		ID: "ch3", DID: "gb-dev-1", DeviceID: "34020000001320000333", Name: "通道1", IsOnline: true,
	}); got != "通道1-门口摄像机" {
		t.Fatalf("got %q", got)
	}
	if got := p.profileName(&ipc.Channel{
		ID: "ch4", DID: "gb-dev-1", Name: "后门", IsOnline: false,
	}); got != "后门-门口摄像机 (离线)" {
		t.Fatalf("got %q", got)
	}
	p.devices["gb-dev-1"] = deviceMeta{name: "门口摄像机", online: false}
	if got := p.profileName(&ipc.Channel{
		ID: "ch5", DID: "gb-dev-1", Name: "侧门", IsOnline: true,
	}); got != "侧门-门口摄像机 (离线)" {
		t.Fatalf("device offline got %q", got)
	}
}

// 同名通道须各占一条 Profile（Token=通道 ID），否则 HA 会合并丢实体。
func TestProfileTokenUniquePerChannel(t *testing.T) {
	p := &StreamProvider{
		devices: map[string]deviceMeta{"dev1": {name: "NVR一楼", online: true}},
		channels: []*ipc.Channel{
			{ID: "ch-online", DID: "dev1", Name: "前门", IsOnline: true},
			{ID: "ch-offline", DID: "dev1", Name: "前门", IsOnline: false},
			{ID: "ch-other", DID: "dev1", Name: "后门", IsOnline: false},
		},
	}
	seen := make(map[string]string)
	for _, ch := range p.channels {
		token := ch.ID
		name := p.profileName(ch)
		if _, ok := seen[token]; ok {
			t.Fatalf("duplicate token %q", token)
		}
		seen[token] = name
	}
	if len(seen) != 3 {
		t.Fatalf("got %d profiles", len(seen))
	}
	if seen["ch-offline"] != "前门-NVR一楼 (离线)" {
		t.Fatalf("offline name %q", seen["ch-offline"])
	}
	if !strings.Contains(seen["ch-other"], "后门") {
		t.Fatalf("other name %q", seen["ch-other"])
	}
}
