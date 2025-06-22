package operation

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestManagerBatchBehavior(t *testing.T) {
	defaultFn := func(ctx context.Context, items Items[string, string]) {
		for _, it := range items {
			it.SetResult(fmt.Sprintf("%s_done", it.Value))
		}
	}

	mgr := NewManager[string](defaultFn, 100*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mgr.Run(ctx)

	id := "user123"
	expected := 3
	var results []string

	for i := 1; i <= expected; i++ {
		val := fmt.Sprintf("val%d", i)
		item, err := mgr.Add(context.Background(), id, expected, val)
		if err != nil {
			t.Fatalf("Add error: %v", err)
		}
		res, err := item.Wait(context.Background())
		if err != nil {
			t.Fatalf("Wait error: %v", err)
		}
		results = append(results, res)
	}

	for i, got := range results {
		want := fmt.Sprintf("val%d_done", i+1)
		if got != want {
			t.Errorf("result[%d] = %q; want %q", i, got, want)
		}
	}
}
