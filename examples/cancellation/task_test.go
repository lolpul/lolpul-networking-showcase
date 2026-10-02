package cancellation_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/lolpul/lolpul-networking-showcase/examples/cancellation"
)

func TestSuccessAndFailure(t *testing.T) {
	failure := errors.New("work failed")
	for _, expected := range []error{nil, failure} {
		task, err := cancellation.Start(context.Background(), func(context.Context) error { return expected })
		if err != nil {
			t.Fatal(err)
		}
		if !errors.Is(task.Wait(), expected) || !errors.Is(task.Stop(), expected) {
			t.Fatal("result changed after join")
		}
	}
}

func TestStopJoinsCleanup(t *testing.T) {
	started, released := make(chan struct{}), make(chan struct{})
	task, _ := cancellation.Start(context.Background(), func(ctx context.Context) error {
		defer close(released)
		close(started)
		<-ctx.Done()
		return ctx.Err()
	})
	<-started
	var group sync.WaitGroup
	for range 16 {
		group.Add(1)
		go func() {
			defer group.Done()
			if !errors.Is(task.Stop(), context.Canceled) {
				t.Error("cancellation lost")
			}
		}()
	}
	group.Wait()
	select {
	case <-released:
	default:
		t.Fatal("join returned before cleanup")
	}
}

func TestParentCancellationReachesWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{})
	task, _ := cancellation.Start(ctx, func(ctx context.Context) error { close(started); <-ctx.Done(); return ctx.Err() })
	<-started
	cancel()
	if !errors.Is(task.Wait(), context.Canceled) {
		t.Fatal("parent cancellation lost")
	}
}

func TestExpiredDeadlineSkipsWork(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Unix(0, 0))
	defer cancel()
	task, _ := cancellation.Start(ctx, func(context.Context) error { t.Error("expired work started"); return nil })
	if !errors.Is(task.Wait(), context.DeadlineExceeded) {
		t.Fatal("deadline lost")
	}
}

func TestDeadlineStopsRunningWorker(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		started := make(chan struct{})
		task, _ := cancellation.Start(ctx, func(ctx context.Context) error {
			close(started)
			<-ctx.Done()
			return ctx.Err()
		})
		<-started
		// When both goroutines block, synctest advances its virtual clock to
		// the deadline. No scheduler-dependent sleep or real-time wait is used.
		if !errors.Is(task.Wait(), context.DeadlineExceeded) {
			t.Fatal("running worker did not observe its deadline")
		}
	})
}

func TestInvalidDependencies(t *testing.T) {
	if _, err := cancellation.Start(nil, func(context.Context) error { return nil }); err == nil {
		t.Fatal("nil context accepted")
	}
	if _, err := cancellation.Start(context.Background(), nil); err == nil {
		t.Fatal("nil work accepted")
	}
}
