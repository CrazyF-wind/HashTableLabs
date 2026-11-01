package hashmap

import (
	"fmt"
	"testing"
)

func TestSetGetUpdateDelete(t *testing.T) {
	m := New()

	m.Set("name", "Alice")
	got, ok := m.Get("name")
	if !ok || got != "Alice" {
		t.Fatalf("Get(name) = (%v, %v), want (Alice, true)", got, ok)
	}

	m.Set("name", "Bob")
	if m.Len() != 1 {
		t.Fatalf("Len() = %d, want 1 after update", m.Len())
	}

	got, ok = m.Get("name")
	if !ok || got != "Bob" {
		t.Fatalf("Get(name) = (%v, %v), want (Bob, true)", got, ok)
	}

	m.Delete("missing")
	if m.Len() != 1 {
		t.Fatalf("Len() = %d, want 1 after deleting missing key", m.Len())
	}

	m.Delete("name")
	if _, ok := m.Get("name"); ok {
		t.Fatal("key should not exist after Delete")
	}
	if m.Len() != 0 {
		t.Fatalf("Len() = %d, want 0 after Delete", m.Len())
	}
}

func TestResizeKeepsEntries(t *testing.T) {
	m := New()
	const count = 1000

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		m.Set(key, i)
	}

	if m.Len() != count {
		t.Fatalf("Len() = %d, want %d", m.Len(), count)
	}

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		got, ok := m.Get(key)
		if !ok || got != i {
			t.Fatalf("Get(%q) = (%v, %v), want (%d, true)", key, got, ok, i)
		}
	}
}

func TestRangeCanStopEarly(t *testing.T) {
	m := New()
	for i := 0; i < 10; i++ {
		m.Set(fmt.Sprintf("key-%d", i), i)
	}

	visited := 0
	m.Range(func(key string, value interface{}) bool {
		visited++
		return false
	})

	if visited != 1 {
		t.Fatalf("visited = %d, want 1", visited)
	}
}
