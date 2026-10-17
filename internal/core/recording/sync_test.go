package recording

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/gowvp/owl/internal/conf"
)

type mockSMS struct {
	streams      map[string]bool
	startedCount atomic.Int32
	stoppedCount atomic.Int32
	startedKeys  []string
	stoppedKeys  []string
}

func (m *mockSMS) StartRecord(app, stream, customPath string, maxSecond int) error {
	m.startedCount.Add(1)
	m.startedKeys = append(m.startedKeys, app+"/"+stream)
	return nil
}

func (m *mockSMS) StopRecord(app, stream string) error {
	m.stoppedCount.Add(1)
	m.stoppedKeys = append(m.stoppedKeys, app+"/"+stream)
	return nil
}

func (m *mockSMS) ListRecordingStreams() (map[string]bool, error) {
	return m.streams, nil
}

type mockIPC struct {
	channels []ChannelInfo
}

func (m *mockIPC) ListOnlineChannels(_ context.Context) ([]ChannelInfo, error) {
	return m.channels, nil
}

type mockPlay struct {
	triggeredCount atomic.Int32
	triggeredIDs   []string
}

func (m *mockPlay) TriggerStream(_ context.Context, info ChannelInfo) error {
	m.triggeredCount.Add(1)
	m.triggeredIDs = append(m.triggeredIDs, info.ID)
	return nil
}

func newTestCore(sms *mockSMS, ipc *mockIPC, play *mockPlay) Core {
	return NewCore(nil,
		WithConfig(&conf.ServerRecording{StorageDir: "/tmp/test"}),
		WithSMSProvider(sms),
		WithIPCProvider(ipc),
		WithPlayProvider(play),
	)
}

// TestSync_StartRecording 应录制+流在但未录制 → startRecord
func TestSync_StartRecording(t *testing.T) {
	sms := &mockSMS{streams: map[string]bool{"rtp/ch1": false}}
	ipc := &mockIPC{channels: []ChannelInfo{
		{ID: "ch1", App: "rtp", Stream: "ch1", RecordMode: "always"},
	}}
	play := &mockPlay{}
	c := newTestCore(sms, ipc, play)

	c.syncRecordingTasks(context.Background())

	if got := sms.startedCount.Load(); got != 1 {
		t.Fatalf("startRecord 调用次数: got %d, want 1", got)
	}
	if got := sms.stoppedCount.Load(); got != 0 {
		t.Fatalf("stopRecord 不应被调用: got %d", got)
	}
	if got := play.triggeredCount.Load(); got != 0 {
		t.Fatalf("triggerStream 不应被调用: got %d", got)
	}
}

// TestSync_StopRecording 不应录制但正在录制 → stopRecord
func TestSync_StopRecording(t *testing.T) {
	sms := &mockSMS{streams: map[string]bool{"rtp/ch1": true}}
	ipc := &mockIPC{channels: []ChannelInfo{
		{ID: "ch1", App: "rtp", Stream: "ch1", RecordMode: "none"},
	}}
	play := &mockPlay{}
	c := newTestCore(sms, ipc, play)

	c.syncRecordingTasks(context.Background())

	if got := sms.stoppedCount.Load(); got != 1 {
		t.Fatalf("stopRecord 调用次数: got %d, want 1", got)
	}
	if got := sms.startedCount.Load(); got != 0 {
		t.Fatalf("startRecord 不应被调用: got %d", got)
	}
}

// TestSync_TriggerStream 应录制但流不在 ZLM → 快照拉流
func TestSync_TriggerStream(t *testing.T) {
	sms := &mockSMS{streams: map[string]bool{}}
	ipc := &mockIPC{channels: []ChannelInfo{
		{ID: "ch1", App: "rtp", Stream: "ch1", RecordMode: "always"},
	}}
	play := &mockPlay{}
	c := newTestCore(sms, ipc, play)

	c.syncRecordingTasks(context.Background())

	if got := play.triggeredCount.Load(); got != 1 {
		t.Fatalf("triggerStream 调用次数: got %d, want 1", got)
	}
	if got := sms.startedCount.Load(); got != 0 {
		t.Fatalf("startRecord 不应被调用: got %d", got)
	}
}

// TestSync_Skip 应录制且已在录制 → 跳过
func TestSync_Skip(t *testing.T) {
	sms := &mockSMS{streams: map[string]bool{"rtp/ch1": true}}
	ipc := &mockIPC{channels: []ChannelInfo{
		{ID: "ch1", App: "rtp", Stream: "ch1", RecordMode: "always"},
	}}
	play := &mockPlay{}
	c := newTestCore(sms, ipc, play)

	c.syncRecordingTasks(context.Background())

	if got := sms.startedCount.Load(); got != 0 {
		t.Fatalf("startRecord 不应被调用: got %d", got)
	}
	if got := sms.stoppedCount.Load(); got != 0 {
		t.Fatalf("stopRecord 不应被调用: got %d", got)
	}
	if got := play.triggeredCount.Load(); got != 0 {
		t.Fatalf("triggerStream 不应被调用: got %d", got)
	}
}

// TestSync_Mixed 混合场景：多通道同时存在不同状态
func TestSync_Mixed(t *testing.T) {
	sms := &mockSMS{streams: map[string]bool{
		"rtp/ch1": true,  // 应录+在录 → skip
		"rtp/ch2": false, // 应录+流在但未录 → start
		"rtp/ch3": true,  // 不应录+在录 → stop
	}}
	ipc := &mockIPC{channels: []ChannelInfo{
		{ID: "ch1", App: "rtp", Stream: "ch1", RecordMode: "always"},
		{ID: "ch2", App: "rtp", Stream: "ch2", RecordMode: "ai"},
		{ID: "ch3", App: "rtp", Stream: "ch3", RecordMode: "none"},
		{ID: "ch4", App: "rtp", Stream: "ch4", RecordMode: "always"},
	}}
	play := &mockPlay{}
	c := newTestCore(sms, ipc, play)

	c.syncRecordingTasks(context.Background())

	if got := sms.startedCount.Load(); got != 1 {
		t.Fatalf("startRecord 调用次数: got %d, want 1 (ch2)", got)
	}
	if got := sms.stoppedCount.Load(); got != 1 {
		t.Fatalf("stopRecord 调用次数: got %d, want 1 (ch3)", got)
	}
	if got := play.triggeredCount.Load(); got != 1 {
		t.Fatalf("triggerStream 调用次数: got %d, want 1 (ch4)", got)
	}
}

// TestSync_EmptyRecordMode 空串 RecordMode 视为 always
func TestSync_EmptyRecordMode(t *testing.T) {
	sms := &mockSMS{streams: map[string]bool{"rtp/ch1": false}}
	ipc := &mockIPC{channels: []ChannelInfo{
		{ID: "ch1", App: "rtp", Stream: "ch1", RecordMode: ""},
	}}
	play := &mockPlay{}
	c := newTestCore(sms, ipc, play)

	c.syncRecordingTasks(context.Background())

	if got := sms.startedCount.Load(); got != 1 {
		t.Fatalf("空串应视为 always, startRecord 调用次数: got %d, want 1", got)
	}
}
