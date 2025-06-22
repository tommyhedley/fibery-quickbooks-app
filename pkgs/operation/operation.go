package operation

import (
	"context"
	"fmt"
	"time"
)

const NoTimeout time.Duration = 0

type Func[T, R any] func(context.Context, Items[T, R])

type Option[T, R any] func(*Operation[T, R])

func WithOpFunc[T, R any](opFn Func[T, R]) Option[T, R] {
	return func(op *Operation[T, R]) {
		op.Fn = opFn
	}
}

func WithTimeout[T, R any](timeout time.Duration) Option[T, R] {
	return func(op *Operation[T, R]) {
		op.timeout = timeout
	}
}

type Operation[T, R any] struct {
	expected int
	Fn       Func[T, R]
	timeout  time.Duration
	input    chan *Item[T, R]
	done     chan struct{}
}

func newOperation[T, R any](expected int, opFn Func[T, R], timeout time.Duration) (*Operation[T, R], error) {
	if expected < 1 {
		return nil, fmt.Errorf("expected is less than 1")
	}

	if opFn == nil {
		return nil, fmt.Errorf("opFn is nil")
	}

	if timeout < 0 {
		return nil, fmt.Errorf("timeout is less than 0")
	}

	return &Operation[T, R]{
		expected: expected,
		Fn:       opFn,
		timeout:  timeout,
		input:    make(chan *Item[T, R]),
		done:     make(chan struct{}),
	}, nil
}

func (op *Operation[T, R]) run(ctx context.Context) {
	defer close(op.done)

	var (
		t *time.Timer
		c <-chan time.Time
	)

	defer func() {
		if t != nil {
			t.Stop()
		}
	}()

	items := make(Items[T, R], 0, op.expected)

	for {
		var run, done, timedOut bool

		select {
		case item := <-op.input:
			items = append(items, item)
			if len(items) == op.expected {
				run = true
				done = true
			}
		case <-c:
			if len(items) > 0 {
				run = true
				timedOut = true
			} else {
				done = true
			}
		case <-ctx.Done():
			if len(items) > 0 {
				run = true
			}
			done = true
		}

		if run {
			op.Fn(ctx, items)

			if timedOut && len(op.input) == 0 {
				break
			}

			c = nil
			items = items[:0]
		}

		if done {
			break
		}

		if !run && c == nil && op.timeout != NoTimeout {
			if t == nil {
				t = time.NewTimer(op.timeout)
			} else {
				t.Reset(op.timeout)
			}
			c = t.C
		}
	}
}
