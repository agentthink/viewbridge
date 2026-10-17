package ffwork

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ixugo/goddd/pkg/queue"
)

type (
	Config struct {
		Width, Height int
		FPS           int
		RTSPURL       string
		Transport     string
		UseWallClock  bool
		HWAccel       string
		OnFrame       func(frame *FrameData)
		Name          string
	}
	FrameData struct {
		FrameNum  uint64
		Timestamp time.Time
		Data      []byte
	}
	FrameCapture struct {
		Name                  string
		config                Config
		frameSize             int
		FrameCh               chan *FrameData
		errCh                 chan error
		ctx                   context.Context
		cancel                context.CancelFunc
		m                     sync.Mutex
		started               bool
		cmd                   *exec.Cmd
		lastFrame             time.Time
		wg                    sync.WaitGroup
		ffmpegLog             *queue.CirQueue[string]
		frameCount, skipCount uint64
		OnFrame               func(frame *FrameData)
	}
	Stats struct {
		Name                  string
		FrameCount, SkipCount uint64
		LastFrame             time.Time
		FrameSize             int
		IsRunning             bool
	}
)

func NewFrameCapture(cfg Config) (*FrameCapture, error) {
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, fmt.Errorf("invalid resolution: %dx%d", cfg.Width, cfg.Height)
	}
	if cfg.FPS <= 0 {
		return nil, fmt.Errorf("invalid fps: %d", cfg.FPS)
	}
	if cfg.RTSPURL == "" {
		return nil, fmt.Errorf("resp url is required")
	}
	if cfg.Transport == "" {
		cfg.Transport = "tcp"
	}
	frameSize := cfg.Width * cfg.Height * 3 / 2
	ctx, cancel := context.WithCancel(context.Background())
	return &FrameCapture{
		config:    cfg,
		frameSize: frameSize,
		FrameCh:   make(chan *FrameData, 10),
		errCh:     make(chan error, 1),
		ctx:       ctx,
		cancel:    cancel,
		ffmpegLog: queue.NewCirQueue[string](100),
		OnFrame:   cfg.OnFrame,
	}, nil
}

func (fc *FrameCapture) FrameSize() int {
	return fc.frameSize
}

func (fc *FrameCapture) buildFFmpegArgs() []string {
	args := []string{
		"-hide_banner",
		"-loglevel", "warning",
		"-threads", "2",
	}
	args = append(args, "-user_agent", "FFmpeg GoWVP")
	args = append(args, "-avoid_negative_ts", "make_zero",
		"-fflags", "+genpts+discardcorrupt",
		"-rtsp_transport", fc.config.Transport,
		"-timeout", "10000000",
	)
	if fc.config.UseWallClock {
		args = append(args, "-use_wallclock_as_timestamps", "1")
	}
	if fc.config.HWAccel != "" {
		args = append(args, "-hwaccel", fc.config.HWAccel)
	}
	args = append(args, "-i", fc.config.RTSPURL)

	args = append(args,
		"-f", "rawvideo",
		"-pix_fmt", "yuv420p",
		"-r", strconv.Itoa(fc.config.FPS),
		"-vf", fmt.Sprintf("fps=%d,scale=%d:%d", fc.config.FPS, fc.config.Width, fc.config.Height),
		"pipe:1",
	)
	return args
}

func (fc *FrameCapture) Start() error {
	fc.m.Lock()
	defer fc.m.Unlock()
	if fc.started {
		return fmt.Errorf("frame capture already started")
	}

	args := fc.buildFFmpegArgs()
	fc.cmd = exec.CommandContext(fc.ctx, "ffmpeg", args...)
	stdout, err := fc.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderr, err := fc.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}
	if err := fc.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}
	fc.started = true
	fc.lastFrame = time.Now()

	fc.wg.Go(func() { fc.captureLoop(stdout) })
	fc.wg.Go(func() { fc.readStderr(stderr) })
	return nil
}

// captureLoop 从 ffmpeg 的 stdout 读取原始视频帧数据
// ffmpeg 输出的是固定大小的 YUV420P 格式帧，需要按帧大小读取
func (fc *FrameCapture) captureLoop(stdout io.Reader) {
	defer close(fc.FrameCh)

	reader := bufio.NewReaderSize(stdout, fc.frameSize*10)
	for {
		select {
		case <-fc.ctx.Done():
			return
		default:
		}

		frameBytes := make([]byte, fc.frameSize)
		n, err := io.ReadFull(reader, frameBytes)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				select {
				case fc.errCh <- fmt.Errorf("ffmpeg stream ended: %w", err):
				default:
				}
				return
			}
			select {
			case fc.errCh <- fmt.Errorf("failed to read frame: %w", err):
			default:
				return
			}
		}
		if n != fc.frameSize {
			select {
			case fc.errCh <- fmt.Errorf("incomplete frame: %d != %d", n, fc.frameSize):
			default:
			}
			return
		}

		frameNum := atomic.AddUint64(&fc.frameCount, 1)
		now := time.Now()
		fc.m.Lock()
		fc.lastFrame = now
		fc.m.Unlock()

		frame := FrameData{
			FrameNum:  frameNum,
			Timestamp: now,
			Data:      frameBytes,
		}
		if fc.OnFrame != nil {
			fc.OnFrame(&frame)
		}

		select {
		case fc.FrameCh <- &frame:
		case <-fc.ctx.Done():
			return
		default:
			atomic.AddUint64(&fc.skipCount, 1)
		}
	}
}

// readStderr 读取 ffmpeg 的 stderr 输出用于日志记录
// ffmpeg 的警告和错误信息都会输出到 stderr
func (fc *FrameCapture) readStderr(stderr io.Reader) {
	scan := bufio.NewScanner(stderr)
	for scan.Scan() {
		line := scan.Text()
		fc.ffmpegLog.Push(line)
	}
}

func (fc *FrameCapture) Frames() <-chan *FrameData {
	return fc.FrameCh
}

func (fc *FrameCapture) Error() <-chan error {
	return fc.errCh
}

func (fc *FrameCapture) Log() []string {
	return fc.ffmpegLog.Range()
}

func (fc *FrameCapture) GetFrame(timeout time.Duration) (*FrameData, error) {
	select {
	case frame, ok := <-fc.FrameCh:
		if !ok {
			return nil, fmt.Errorf("frame channel closed")
		}
		return frame, nil
	case err := <-fc.errCh:
		return nil, err
	case <-fc.ctx.Done():
		return nil, fc.ctx.Err()
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout")
	}
}

func (fc *FrameCapture) Stop() error {
	fc.m.Lock()
	if !fc.started {
		fc.m.Unlock()
		return nil
	}

	fc.m.Unlock()
	if cancel := fc.cancel; cancel != nil {
		cancel()
	}
	fc.wg.Wait()

	if fc.cmd != nil && fc.cmd.Process != nil {
		done := make(chan error, 1)
		defer close(done)
		go func() {
			done <- fc.cmd.Wait()
		}()

		select {
		case <-time.After(5 * time.Second):
			if err := fc.cmd.Process.Kill(); err != nil {
				return fmt.Errorf("failed to kill ffmpeg: %w", err)
			}
			<-done
		case <-done:
		}
	}
	return nil
}

func (fc *FrameCapture) GetStats() Stats {
	fc.m.Lock()
	defer fc.m.Unlock()
	return Stats{
		Name:       fc.config.Name,
		FrameCount: atomic.LoadUint64(&fc.frameCount),
		SkipCount:  atomic.LoadUint64(&fc.skipCount),
		LastFrame:  fc.lastFrame,
		FrameSize:  fc.frameSize,
		IsRunning:  fc.started,
	}
}
