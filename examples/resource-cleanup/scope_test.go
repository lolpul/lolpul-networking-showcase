package cleanup_test

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	cleanup "github.com/lolpul/lolpul-networking-showcase/examples/resource-cleanup"
)

type resource struct {
	name    string
	order   *[]string
	closes  int
	failure error
}

func (r *resource) Close() error { r.closes++; *r.order = append(*r.order, r.name); return r.failure }

func TestSuccessAndEveryFailureBoundary(t *testing.T) {
	failure := errors.New("acquisition or use failure")
	for _, stage := range []string{"success", "first", "first-partial", "second", "second-partial", "use"} {
		t.Run(stage, func(t *testing.T) {
			order := []string{}
			a, b := &resource{name: "first", order: &order}, &resource{name: "second", order: &order}
			calls := 0
			err := cleanup.WithPair(context.Background(), func(context.Context) (io.Closer, error) {
				if stage == "first" {
					return nil, failure
				}
				if stage == "first-partial" {
					return a, failure
				}
				return a, nil
			}, func(context.Context) (io.Closer, error) {
				if stage == "second" {
					return nil, failure
				}
				if stage == "second-partial" {
					return b, failure
				}
				return b, nil
			}, func(context.Context, io.Closer, io.Closer) error {
				calls++
				if stage == "use" {
					return failure
				}
				return nil
			})
			wantOrder, wantCalls := []string{}, 0
			if stage != "first" {
				wantOrder = []string{"first"}
			}
			if stage == "success" || stage == "use" || stage == "second-partial" {
				wantOrder = []string{"second", "first"}
			}
			if stage == "success" || stage == "use" {
				wantCalls = 1
			}
			if !reflect.DeepEqual(order, wantOrder) || calls != wantCalls || a.closes > 1 || b.closes > 1 {
				t.Fatalf("ownership: %v, use=%d", order, calls)
			}
			if stage == "success" && err != nil {
				t.Fatal(err)
			}
			if stage != "success" && !errors.Is(err, failure) {
				t.Fatal("failure lost")
			}
		})
	}
}

func TestAllCleanupErrorsSurvive(t *testing.T) {
	order := []string{}
	operation, closeA, closeB := errors.New("operation"), errors.New("close first"), errors.New("close second")
	a, b := &resource{failure: closeA, order: &order}, &resource{failure: closeB, order: &order}
	err := cleanup.WithPair(context.Background(), func(context.Context) (io.Closer, error) { return a, nil }, func(context.Context) (io.Closer, error) { return b, nil }, func(context.Context, io.Closer, io.Closer) error { return operation })
	for _, expected := range []error{operation, closeA, closeB} {
		if !errors.Is(err, expected) {
			t.Fatal("joined error lost")
		}
	}
}

func TestCancellationDuringAcquisitionReleasesOwnership(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	order := []string{}
	a := &resource{name: "first", order: &order}
	err := cleanup.WithPair(ctx, func(context.Context) (io.Closer, error) { cancel(); return a, nil }, func(context.Context) (io.Closer, error) {
		t.Error("second acquisition after cancellation")
		return nil, nil
	}, func(context.Context, io.Closer, io.Closer) error { t.Error("use after cancellation"); return nil })
	if !errors.Is(err, context.Canceled) || a.closes != 1 {
		t.Fatal("cancelled acquisition leaked")
	}
}

func TestMissingResourceRejected(t *testing.T) {
	for _, missingFirst := range []bool{true, false} {
		order := []string{}
		a := &resource{order: &order}
		err := cleanup.WithPair(context.Background(), func(context.Context) (io.Closer, error) {
			if missingFirst {
				return nil, nil
			}
			return a, nil
		}, func(context.Context) (io.Closer, error) { return nil, nil }, func(context.Context, io.Closer, io.Closer) error { t.Error("use with missing resource"); return nil })
		if err == nil {
			t.Fatal("missing resource accepted")
		}
		if !missingFirst && a.closes != 1 {
			t.Fatal("first resource leaked")
		}
	}
}
