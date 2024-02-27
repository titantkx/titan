package testutil

import (
	"io"
)

// Streamer allows to write to and read from it from different goroutines
type Streamer struct {
	ch chan []byte
}

func NewStreamer() *Streamer {
	return &Streamer{ch: make(chan []byte)}
}

func (s *Streamer) Write(p []byte) (int, error) {
	s.ch <- p
	return len(p), nil
}

func (s *Streamer) Read(p []byte) (int, error) {
	buf, ok := <-s.ch
	if !ok {
		return 0, io.EOF
	}
	return copy(p, buf), nil
}

func (s *Streamer) Close() error {
	close(s.ch)
	return nil
}
