package operation

import (
	"context"
	"testing"
	"time"
)

func TestOperation_ItemContextCancellation(t *testing.T) {
	called := false
	op, err := newOperation(2, func(ctx context.Context, items Items[int, int]) {
		called = true
		for _, it := range items {
			it.SetResult(it.Value)
		}
	}, 100*time.Millisecond)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	go op.run(context.Background())

	// create an item whose ctx we’ll cancel before sending
	cancelCtx, cancel := context.WithCancel(context.Background())
	item := newItem[int, int](42)
	cancel() // cancel immediately

	// try to add via the unexported channel to simulate Add
	select {
	case op.input <- item:
	case <-time.After(10 * time.Millisecond):
		t.Fatal("could not enqueue item")
	}

	// Wait with the cancelled context
	if _, err := item.Wait(cancelCtx); err == nil {
		t.Error("expected error from cancelled context, got nil")
	}
	if called {
		t.Error("batch function should not have been called on cancelled item")
	}
}

func TestOperation_RootContextCancellation(t *testing.T) {
	// use a channel to detect fn invocation
	fnStarted := make(chan struct{})
	op, err := newOperation(3, func(ctx context.Context, items Items[int, int]) {
		close(fnStarted)
		for _, it := range items {
			it.SetResult(it.Value)
		}
	}, 1*time.Second) // long timer so batch won’t auto-fire

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rootCtx, cancelRoot := context.WithCancel(context.Background())
	go op.run(rootCtx)

	// send one item
	item1 := newItem[int, int](10)
	op.input <- item1

	// now cancel the root context before we hit expected=3
	cancelRoot()

	// item1.Wait should return ctx.Err()
	if _, err := item1.Wait(rootCtx); err == nil {
		t.Error("expected root context cancellation error, got nil")
	}

	// verify batch func never started
	select {
	case <-fnStarted:
		t.Error("batch function should not have been invoked after root cancel")
	default:
	}
}

func TestManager_ItemContextCancel(t *testing.T) {
	defaultFn := func(ctx context.Context, items Items[string, string]) {}
	mgr := NewManager[string](defaultFn, 50*time.Millisecond)
	appCtx, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()
	mgr.Run(appCtx)

	id := "testID"
	reqCtx, cancelReq := context.WithCancel(context.Background())
	item, err := mgr.Add(reqCtx, id, 1, "foo")
	if err != nil {
		t.Fatalf("unexpected Add error: %v", err)
	}
	cancelReq()
	if _, err := item.Wait(reqCtx); err == nil {
		t.Error("expected error from cancelled request context, got nil")
	}
}

func TestManager_AppContextCancel(t *testing.T) {
	called := false
	defaultFn := func(ctx context.Context, items Items[int, int]) {
		called = true
		for _, it := range items {
			it.SetResult(it.Value)
		}
	}

	mgr := NewManager[string](defaultFn, 1*time.Second)
	appCtx, cancelApp := context.WithCancel(context.Background())
	mgr.Run(appCtx)

	// queue up one item (won’t fire until timeout or appCtx.Done)
	item, err := mgr.Add(context.Background(), "id", 5, 123)
	if err != nil {
		t.Fatalf("unexpected Add error: %v", err)
	}

	// shut down the app — now Run will stop accepting new reqs
	cancelApp()

	// • Expect no error, because we're waiting with Background
	// • Expect to see the batch fn run and produce our result
	res, err := item.Wait(context.Background())
	if err != nil {
		t.Errorf("expected no error after app shutdown, got %v", err)
	}
	if res != 123 {
		t.Errorf("expected value 123, got %v", res)
	}
	if !called {
		t.Error("expected batch fn to run after app shutdown")
	}
}
