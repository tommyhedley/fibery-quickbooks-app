package cache

import (
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	t.Parallel()
	c := NewCache[string, int](time.Second)
	c.Set("foo", 42)
	v, ok := c.Get("foo")
	if !ok {
		t.Fatal("expected key 'foo' to be found")
	}
	if v != 42 {
		t.Fatalf("expected value 42, got %v", v)
	}
}

func TestExpiration(t *testing.T) {
	t.Parallel()
	c := NewCache[string, int](10 * time.Millisecond)
	c.Set("bar", 100)

	time.Sleep(15 * time.Millisecond)
	if _, ok := c.Get("bar"); ok {
		t.Fatal("expected key 'bar' to have expired and been deleted")
	}
}

func TestCleanup(t *testing.T) {
	t.Parallel()
	c := NewCache[string, int](5 * time.Millisecond)
	c.SetWithTTL("a", 1, 5*time.Millisecond)
	c.SetWithTTL("b", 2, 100*time.Millisecond)

	time.Sleep(10 * time.Millisecond)
	c.Cleanup()

	if _, ok := c.Get("a"); ok {
		t.Error("expected 'a' to be cleaned up")
	}
	if v, ok := c.Get("b"); !ok || v != 2 {
		t.Errorf("expected 'b'=2 to remain, got %v (found=%v)", v, ok)
	}
}
