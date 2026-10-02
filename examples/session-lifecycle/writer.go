// Package lifecycle demonstrates one owner for a serial, closable writer.
package lifecycle

import (
	"errors"
	"io"
	"sync"
)

var ErrFinished = errors.New("writer session has finished")

type Phase uint8

const (
	Ready Phase = iota
	Finished
)

// WriterSession takes ownership of a non-nil writer on construction. The writer
// need not be concurrency-safe: every write and close is serialized here.
// Do not use or close the underlying writer after transferring ownership.
type WriterSession struct {
	mu       sync.Mutex
	writer   io.WriteCloser
	phase    Phase
	closeErr error
}

func New(writer io.WriteCloser) (*WriterSession, error) {
	if writer == nil {
		return nil, errors.New("writer is required")
	}
	return &WriterSession{writer: writer}, nil
}

func (s *WriterSession) Phase() Phase {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.phase
}

func (s *WriterSession) Write(data []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.phase == Finished {
		return 0, ErrFinished
	}
	n, err := s.writer.Write(data)
	if n < len(data) && err == nil {
		err = io.ErrShortWrite
	}
	return n, err
}

// Close waits for any current Write, then attempts cleanup exactly once. A close
// failure is retained, not retried: the writer's recovery semantics are unknown.
// Writer methods must terminate and must not call back into this session.
func (s *WriterSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.phase == Ready {
		s.phase = Finished
		s.closeErr = s.writer.Close()
		s.writer = nil
	}
	return s.closeErr
}

var _ io.WriteCloser = (*WriterSession)(nil)
