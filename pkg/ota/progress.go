package ota

import (
	"io"
	"sync/atomic"
	"time"
)

type ProgressReader struct {
	Total   int64
	Current atomic.Int64
	io.Reader
	OnProgress func(current, total int64)
	quit       chan struct{}
}

func NewProgressReader(total int64, reader io.Reader, onProgress func(current, total int64)) *ProgressReader {
	p := ProgressReader{
		Total:      total,
		Reader:     reader,
		OnProgress: onProgress,
		quit:       make(chan struct{}, 1),
	}
	if onProgress != nil {
		go p.Start()
	}

	return &p
}

func (p *ProgressReader) Close() {
	close(p.quit)
}

func (p *ProgressReader) Start() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.OnProgress(p.Current.Load(), p.Total)
		case <-p.quit:
			p.OnProgress(p.Current.Load(), p.Total)
			return
		}
	}
}

func (p *ProgressReader) Read(b []byte) (int, error) {
	n, err := p.Reader.Read(b)
	p.Current.Add(int64(n))
	return n, err
}
