package cache

import (
	"testing"
	"time"
)

func TestMemoryCache_SetAndGet(t *testing.T) {
	c := NewMemoryCache(1 * time.Minute)
	defer c.Close()

	c.Set("key1", []byte("value1"), 5*time.Minute)

	data, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if string(data) != "value1" {
		t.Errorf("expected value1, got %s", string(data))
	}
}

func TestMemoryCache_Miss(t *testing.T) {
	c := NewMemoryCache(1 * time.Minute)
	defer c.Close()

	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected cache miss")
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	c := NewMemoryCache(50 * time.Millisecond)
	defer c.Close()

	c.Set("expiring", []byte("data"), 100*time.Millisecond)

	data, ok := c.Get("expiring")
	if !ok {
		t.Fatal("expected cache hit before expiration")
	}
	if string(data) != "data" {
		t.Errorf("expected data, got %s", string(data))
	}

	time.Sleep(150 * time.Millisecond)

	_, ok = c.Get("expiring")
	if ok {
		t.Error("expected cache miss after expiration")
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	c := NewMemoryCache(1 * time.Minute)
	defer c.Close()

	c.Set("del", []byte("val"), 5*time.Minute)
	c.Delete("del")

	_, ok := c.Get("del")
	if ok {
		t.Error("expected cache miss after delete")
	}
}

func TestMemoryCache_Flush(t *testing.T) {
	c := NewMemoryCache(1 * time.Minute)
	defer c.Close()

	c.Set("a", []byte("1"), 5*time.Minute)
	c.Set("b", []byte("2"), 5*time.Minute)

	c.Flush()

	_, ok1 := c.Get("a")
	_, ok2 := c.Get("b")
	if ok1 || ok2 {
		t.Error("expected all entries flushed")
	}
}

func TestMemoryCache_Overwrite(t *testing.T) {
	c := NewMemoryCache(1 * time.Minute)
	defer c.Close()

	c.Set("key", []byte("v1"), 5*time.Minute)
	c.Set("key", []byte("v2"), 5*time.Minute)

	data, ok := c.Get("key")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if string(data) != "v2" {
		t.Errorf("expected v2, got %s", string(data))
	}
}
