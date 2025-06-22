package operation

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Manager[I comparable, T, R any] struct {
	mu             sync.Mutex
	ops            map[I]*Operation[T, R]
	defaultFn      Func[T, R]
	defaultTimeout time.Duration
	submit         chan Request[I, T, R]
}

type Request[I comparable, T, R any] struct {
	id       I
	expected int
	item     *Item[T, R]
	options  []Option[T, R]
}

func NewManager[I comparable, T, R any](defaultFn Func[T, R], defaultTimeout time.Duration) *Manager[I, T, R] {
	return &Manager[I, T, R]{
		ops:            make(map[I]*Operation[T, R]),
		defaultFn:      defaultFn,
		defaultTimeout: defaultTimeout,
		submit:         make(chan Request[I, T, R]),
	}
}

func (m *Manager[I, T, R]) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case req := <-m.submit:
				m.mu.Lock()
				op := m.ops[req.id]
				if op != nil {
					select {
					case <-op.done:
						delete(m.ops, req.id)
						op = nil
					default:
					}
				}
				if op == nil {
					var err error
					op, err = newOperation(req.expected, m.defaultFn, m.defaultTimeout)
					if err != nil {
						req.item.SetError(fmt.Errorf("error building operation: %w", err))
					}

					for _, opt := range req.options {
						opt(op)
					}

					if op.Fn == nil {
						req.item.SetError(fmt.Errorf("error building operation: opFn was nil"))
						continue
					}

					if op.timeout < 0 {
						req.item.SetError(fmt.Errorf("error building operation: timeout was less than 0"))
						continue
					}

					m.ops[req.id] = op

					go func(id I, me *Operation[T, R]) {
						me.run(ctx)
						m.mu.Lock()

						if current, ok := m.ops[id]; ok && current == me {
							delete(m.ops, id)
						}

						m.mu.Unlock()
					}(req.id, op)
				}
				m.mu.Unlock()

				op.input <- req.item
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (m *Manager[I, T, R]) Add(ctx context.Context, opID I, expected int, v T, opts ...Option[T, R]) (*Item[T, R], error) {
	item := newItem[T, R](v)
	req := Request[I, T, R]{
		id:       opID,
		expected: expected,
		item:     item,
		options:  opts,
	}

	select {
	case m.submit <- req:
		return item, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
