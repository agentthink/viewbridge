package ipc

import "testing"

func TestDeviceGetName(t *testing.T) {
	d := &Device{Name: "", Ext: DeviceExt{Name: "上报名"}, DeviceID: "34010000001320000001"}
	if got := d.GetName(); got != "上报名" {
		t.Fatalf("got %q", got)
	}
	d2 := &Device{Name: "自定义", Ext: DeviceExt{Name: "上报名"}}
	if got := d2.GetName(); got != "自定义" {
		t.Fatalf("got %q", got)
	}
}

func TestProfileDisplayName(t *testing.T) {
	ch := &Channel{Name: "通道1"}
	if got := ProfileDisplayName(ch, "一楼NVR"); got != "通道1-一楼NVR" {
		t.Fatalf("got %q", got)
	}
}

func TestChannelReachable(t *testing.T) {
	ch := &Channel{IsOnline: true}
	if !ChannelReachable(ch, true) {
		t.Fatal("want reachable")
	}
	if ChannelReachable(ch, false) {
		t.Fatal("device offline should block")
	}
	ch.IsOnline = false
	if ChannelReachable(ch, true) {
		t.Fatal("channel offline should block")
	}
}

func TestONVIFProfileNameRTSPNoOfflineSuffix(t *testing.T) {
	ch := &Channel{ID: "sp1", Name: "拉流1", Type: TypeRTSP, IsOnline: false}
	if got := ONVIFProfileName(ch, "dev", false); got != "拉流1-dev" {
		t.Fatalf("rtsp should not suffix offline, got %q", got)
	}
}
