// Package cleanup demonstrates scoped ownership across partial acquisition.
package cleanup

import (
	"context"
	"errors"
	"io"
)

// Acquire may return a resource together with an error. If a non-nil resource
// is returned, ownership transfers to WithPair even when acquisition fails.
// Return a nil interface, not a typed nil. Close must terminate without callbacks.
type Acquire func(context.Context) (io.Closer, error)

// WithPair owns two temporary resources. Use borrows them only for its duration;
// it must not close or retain them. No goroutines are created by this function.
func WithPair(ctx context.Context, first, second Acquire, use func(context.Context, io.Closer, io.Closer) error) (result error) {
	if ctx == nil || first == nil || second == nil || use == nil {
		return errors.New("context, acquisitions and use are required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	a, err := first(ctx)
	if a != nil {
		defer func() { result = errors.Join(result, a.Close()) }()
	}
	if err != nil {
		return err
	}
	if a == nil {
		return errors.New("first acquisition returned no resource")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	b, err := second(ctx)
	if b != nil {
		defer func() { result = errors.Join(result, b.Close()) }()
	}
	if err != nil {
		return err
	}
	if b == nil {
		return errors.New("second acquisition returned no resource")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return use(ctx, a, b)
}
