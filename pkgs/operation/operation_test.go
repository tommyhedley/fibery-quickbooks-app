package operation

import (
	"context"
	"testing"
	"time"
)

// TestOperationBatchByCount verifies that Operation.run fires Fn
// exactly when expected items are enqueued.
func TestOperationBatchByCount(t *testing.T) {
	expected := 3
	op, err := newOperation(expected, func(ctx context.Context, items Items[int, int]) {
		for _, it := range items {
			it.SetResult(it.Value * 2) // multiply each value by 2
		}
	}, time.Hour)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	go op.run(context.Background())

	// enqueue exactly `expected` items
	items := make([]*Item[int, int], expected)
	for i := 0; i < expected; i++ {
		items[i] = newItem[int, int](i + 1)
		op.input <- items[i]
	}

	// each Wait should return double the input
	for _, it := range items {
		res, err := it.Wait(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := it.Value * 2
		if res != want {
			t.Errorf("got %d; want %d", res, want)
		}
	}
}

// TestOperationBatchByTimeout verifies that when fewer than expected items
// are enqueued, they flush after timeout.
func TestOperationBatchByTimeout(t *testing.T) {
	expected := 5
	timeout := 50 * time.Millisecond
	collected := make([]int, 0, expected)
	op, err := newOperation(expected, func(ctx context.Context, items Items[int, int]) {
		for _, it := range items {
			collected = append(collected, it.Value)
			it.SetResult(it.Value)
		}
	}, timeout)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	go op.run(context.Background())

	// enqueue fewer than expected
	items := []*Item[int, int]{
		newItem[int, int](10),
		newItem[int, int](20),
	}
	for _, it := range items {
		op.input <- it
	}

	// wait for timeout + small buffer
	time.Sleep(timeout + 20*time.Millisecond)

	// each Wait should complete without error
	for _, it := range items {
		if _, err := it.Wait(context.Background()); err != nil {
			t.Fatalf("unexpected error on Wait(): %v", err)
		}
	}

	// collected should match the two values
	if len(collected) != len(items) {
		t.Fatalf("Fn was called with %d items; want %d", len(collected), len(items))
	}
	for i, v := range collected {
		if v != items[i].Value {
			t.Errorf("collected[%d] = %d; want %d", i, v, items[i].Value)
		}
	}
}
