// Package cancellation demonstrates cancellation followed by joining a worker.
package cancellation

import (
	"context"
	"errors"
)

// Work must honor ctx cancellation, release everything it acquires, and return.
// A context cannot forcibly interrupt arbitrary blocking code.
type Work func(ctx context.Context) error

type Task struct {
	cancel context.CancelFunc
	done   chan struct{}
	result error // Written before done closes; readers wait for that publication.
}

// Start borrows the parent context and owns exactly one worker goroutine.
// The caller must eventually Wait or Stop; otherwise its lifetime is unbounded.
func Start(parent context.Context, work Work) (*Task, error) {
	if parent == nil || work == nil {
		return nil, errors.New("context and work are required")
	}
	ctx, cancel := context.WithCancel(parent)
	t := &Task{cancel: cancel, done: make(chan struct{})}
	go func() {
		defer close(t.done)
		defer cancel()
		if err := ctx.Err(); err != nil {
			t.result = err
			return
		}
		t.result = work(ctx)
	}()
	return t, nil
}

// Wait joins the worker. Its result is the worker's result; cancellation arriving
// after successful work does not turn that success into an artificial error.
func (t *Task) Wait() error {
	<-t.done
	return t.result
}

// Stop can be called concurrently or repeatedly. Cancellation alone is not
// proof of cleanup; waiting for done is the ownership boundary.
func (t *Task) Stop() error {
	t.cancel()
	return t.Wait()
}
