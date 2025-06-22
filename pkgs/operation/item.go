package operation

import "context"

type Item[T, R any] struct {
	Value  T
	result R
	err    error
	done   chan struct{}
}

func newItem[T, R any](v T) *Item[T, R] {
	return &Item[T, R]{
		Value: v,
		done:  make(chan struct{}),
	}
}

func (i *Item[T, R]) SetResult(result R) {
	i.result = result
	close(i.done)
}

func (i *Item[T, R]) SetError(err error) {
	i.err = err
	close(i.done)
}

func (i *Item[T, R]) Wait(ctx context.Context) (R, error) {
	var zero R
	select {
	case <-i.done:
		if i.err != nil {
			return zero, i.err
		}
		return i.result, nil
	case <-ctx.Done():
		return zero, ctx.Err()
	}
}

type Items[T, R any] []*Item[T, R]

func (is Items[T, R]) SetError(err error) {
	for _, i := range is {
		i.SetError(err)
	}
}
