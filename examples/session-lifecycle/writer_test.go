package lifecycle_test

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"testing"

	lifecycle "github.com/lolpul/lolpul-networking-showcase/examples/session-lifecycle"
)

type sink struct {
	bytes.Buffer
	closed int
	err    error
	short  bool
}

func (w *sink) Close() error { w.closed++; return w.err }
func (w *sink) Write(p []byte) (int, error) {
	if w.short {
		return 0, nil
	}
	return w.Buffer.Write(p)
}

func TestOwnershipAndRepeatedClose(t *testing.T) {
	failure := errors.New("cleanup failure")
	w := &sink{err: failure}
	s, err := lifecycle.New(w)
	if err != nil {
		t.Fatal(err)
	}
	if s.Phase() != lifecycle.Ready {
		t.Fatal("incorrect initial phase")
	}
	if n, err := s.Write([]byte("sample")); n != 6 || err != nil {
		t.Fatalf("write: %d %v", n, err)
	}
	for range 2 {
		if !errors.Is(s.Close(), failure) {
			t.Fatal("cleanup error lost")
		}
	}
	if w.closed != 1 || w.String() != "sample" || s.Phase() != lifecycle.Finished {
		t.Fatal("ownership violated")
	}
	if _, err := s.Write(nil); !errors.Is(err, lifecycle.ErrFinished) {
		t.Fatal("write after close accepted")
	}
}

func TestInvalidWriterAndShortWrite(t *testing.T) {
	if _, err := lifecycle.New(nil); err == nil {
		t.Fatal("nil writer accepted")
	}
	s, _ := lifecycle.New(&sink{short: true})
	defer s.Close()
	if _, err := s.Write([]byte("sample")); !errors.Is(err, io.ErrShortWrite) {
		t.Fatal("short write lost")
	}
}

type blockingSink struct {
	started chan struct{}
	finish  chan struct{}
	closed  chan struct{}
}

func (w *blockingSink) Write(p []byte) (int, error) { close(w.started); <-w.finish; return len(p), nil }
func (w *blockingSink) Close() error                { close(w.closed); return nil }

func TestCloseWaitsForWrite(t *testing.T) {
	w := &blockingSink{make(chan struct{}), make(chan struct{}), make(chan struct{})}
	s, _ := lifecycle.New(w)
	written := make(chan struct{})
	go func() { defer close(written); _, _ = s.Write([]byte("sample")) }()
	<-w.started
	closed := make(chan struct{})
	go func() { defer close(closed); _ = s.Close() }()
	select {
	case <-w.closed:
		t.Fatal("closed during write")
	default:
	}
	close(w.finish)
	<-written
	<-closed
	select {
	case <-w.closed:
	default:
		t.Fatal("cleanup missing")
	}
}

func TestConcurrentWriteObserveAndClose(t *testing.T) {
	w := &sink{}
	s, _ := lifecycle.New(w)
	var group sync.WaitGroup
	for range 32 {
		group.Add(1)
		go func() {
			defer group.Done()
			_, err := s.Write([]byte("x"))
			if err != nil && !errors.Is(err, lifecycle.ErrFinished) {
				t.Error(err)
			}
			_ = s.Phase()
			if err := s.Close(); err != nil {
				t.Error(err)
			}
		}()
	}
	group.Wait()
	if w.closed != 1 {
		t.Fatal("cleanup repeated")
	}
}
