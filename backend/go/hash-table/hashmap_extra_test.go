package hashmap

import (
	"fmt"
	"math/rand"
	"testing"
)

// pickCollidingKeys 返回 n 个哈希到同一个桶的 key。
// 这是白盒测试：直接用 index 判断，不依赖对 FNV 分布的猜测。
func pickCollidingKeys(t *testing.T, m *HashMap, n int) []string {
	t.Helper()

	groups := make(map[int][]string)
	for i := 0; i < 1000000; i++ {
		key := fmt.Sprintf("collide-%d", i)
		idx := m.index(key)
		groups[idx] = append(groups[idx], key)
		if len(groups[idx]) >= n {
			return groups[idx][:n]
		}
	}
	t.Fatalf("failed to find %d keys in the same bucket", n)
	return nil
}

// chainKeys 从桶里读出链表上的 key，按从表头到表尾的顺序返回。
func chainKeys(m *HashMap, idx int) []string {
	var keys []string
	for e := m.buckets[idx]; e != nil; e = e.next {
		keys = append(keys, e.key)
	}
	return keys
}

func TestHashStringMatchesFNV1a(t *testing.T) {
	cases := []struct {
		input string
		want  uint32
	}{
		{"", 0x811c9dc5},
		{"a", 0xe40c292c},
		{"foobar", 0xbf9cf968},
	}
	for _, c := range cases {
		if got := hashString(c.input); got != c.want {
			t.Errorf("hashString(%q) = %#x, want %#x", c.input, got, c.want)
		}
	}
}

func TestUpdateDoesNotGrow(t *testing.T) {
	m := New()
	keys := pickCollidingKeys(t, m, 6) // 6/8 == 0.75，未超过阈值

	for _, k := range keys {
		m.Set(k, 1)
	}
	if len(m.buckets) != 8 {
		t.Fatalf("buckets = %d, want 8 before exceeding loadFactor", len(m.buckets))
	}

	// 反复更新已有 key：不应增加 size，也不应触发扩容
	for i := 0; i < 100; i++ {
		m.Set(keys[0], i)
	}
	if len(m.buckets) != 8 {
		t.Fatalf("buckets = %d, want 8 after updates", len(m.buckets))
	}
	if m.Len() != 6 {
		t.Fatalf("Len() = %d, want 6", m.Len())
	}
	if v, ok := m.Get(keys[0]); !ok || v != 99 {
		t.Fatalf("Get(%q) = (%v, %v), want (99, true)", keys[0], v, ok)
	}

	// 第 7 个新 key 使 7/8 > 0.75，才触发扩容
	m.Set("one-more-key", 0)
	if len(m.buckets) != 16 {
		t.Fatalf("buckets = %d, want 16 after exceeding loadFactor", len(m.buckets))
	}
}

func TestDeleteHeadMiddleTail(t *testing.T) {
	m := New()
	keys := pickCollidingKeys(t, m, 3) // keys[0] -> keys[1] -> keys[2]

	for _, k := range keys {
		m.Set(k, k)
	}

	idx := m.index(keys[0])
	// 头插：后插入的排在链首
	if got := chainKeys(m, idx); len(got) != 3 || got[0] != keys[2] || got[1] != keys[1] || got[2] != keys[0] {
		t.Fatalf("chain = %v, want [%s %s %s]", got, keys[2], keys[1], keys[0])
	}

	// 删链表中间节点：前驱的 next 必须指向后继
	m.Delete(keys[1])
	if got := chainKeys(m, idx); len(got) != 2 || got[0] != keys[2] || got[1] != keys[0] {
		t.Fatalf("chain after deleting middle = %v, want [%s %s]", got, keys[2], keys[0])
	}
	if _, ok := m.Get(keys[1]); ok {
		t.Fatal("middle key still found after Delete")
	}
	if m.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", m.Len())
	}

	// 删链表头节点：桶头指针必须后移
	m.Delete(keys[2])
	if got := chainKeys(m, idx); len(got) != 1 || got[0] != keys[0] {
		t.Fatalf("chain after deleting head = %v, want [%s]", got, keys[0])
	}

	// 删链表尾节点（也是唯一节点）：桶必须变空
	m.Delete(keys[0])
	if got := chainKeys(m, idx); len(got) != 0 {
		t.Fatalf("chain after deleting tail = %v, want empty", got)
	}
	if m.Len() != 0 {
		t.Fatalf("Len() = %d, want 0", m.Len())
	}
}

func TestResizeKeepsCollidingKeys(t *testing.T) {
	m := New()
	keys := pickCollidingKeys(t, m, 12) // 在 8 桶下全部同桶，插入过程会触发扩容

	for i, k := range keys {
		m.Set(k, i)
	}
	if len(m.buckets) != 16 {
		t.Fatalf("buckets = %d, want 16", len(m.buckets))
	}
	if m.Len() != len(keys) {
		t.Fatalf("Len() = %d, want %d", m.Len(), len(keys))
	}

	for i, k := range keys {
		got, ok := m.Get(k)
		if !ok || got != i {
			t.Fatalf("Get(%q) = (%v, %v), want (%d, true)", k, got, ok, i)
		}
	}
}

func TestResizeSplitsIntoAtMostTwoBuckets(t *testing.T) {
	m := New()
	oldLen := len(m.buckets) // 8
	oldIdx := m.index("split-probe")

	// 造 6 个落在 oldIdx 的 key：6/8 == 0.75，不会触发扩容
	var keys []string
	for i := 0; len(keys) < 6; i++ {
		k := fmt.Sprintf("split-%d", i)
		if m.index(k) == oldIdx {
			keys = append(keys, k)
		}
	}
	for _, k := range keys {
		m.Set(k, k)
	}
	if len(m.buckets) != oldLen {
		t.Fatalf("unexpected resize before trigger: buckets = %d, want %d", len(m.buckets), oldLen)
	}

	// 第 7 个 key 使 7/8 > 0.75，精确触发一次扩容
	m.Set("trigger-resize", 0)

	if len(m.buckets) != oldLen*2 {
		t.Fatalf("buckets = %d, want %d", len(m.buckets), oldLen*2)
	}
	for _, k := range keys {
		idx := m.index(k)
		if idx != oldIdx && idx != oldIdx+oldLen {
			t.Fatalf("key %q moved to bucket %d, want %d or %d", k, idx, oldIdx, oldIdx+oldLen)
		}
	}
}

func TestRangeVisitsEveryEntryOnce(t *testing.T) {
	m := New()
	const count = 500
	for i := 0; i < count; i++ {
		m.Set(fmt.Sprintf("key-%d", i), i)
	}

	seen := make(map[string]interface{})
	m.Range(func(key string, value interface{}) bool {
		if _, dup := seen[key]; dup {
			t.Fatalf("key %q visited twice", key)
		}
		seen[key] = value
		return true
	})

	if len(seen) != count {
		t.Fatalf("visited %d entries, want %d", len(seen), count)
	}
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		if v, ok := seen[key]; !ok || v != i {
			t.Fatalf("visited[%q] = (%v, %v), want (%d, true)", key, v, ok, i)
		}
	}
}

func TestGetMissingReturnsFalse(t *testing.T) {
	m := New()
	m.Set("exists", 1)

	if v, ok := m.Get("missing"); ok || v != nil {
		t.Fatalf("Get(missing) = (%v, %v), want (nil, false)", v, ok)
	}
}

// TestDifferentialAgainstNativeMap 用随机操作序列与 Go 原生 map 做差分验证。
func TestDifferentialAgainstNativeMap(t *testing.T) {
	m := New()
	want := make(map[string]int)
	r := rand.New(rand.NewSource(20261030))

	for i := 0; i < 20000; i++ {
		key := fmt.Sprintf("k-%d", r.Intn(500))

		if r.Intn(10) < 3 { // 30% 删除，70% 写入/更新
			m.Delete(key)
			delete(want, key)
		} else {
			v := r.Int()
			m.Set(key, v)
			want[key] = v
		}

		if m.Len() != len(want) {
			t.Fatalf("step %d: Len() = %d, want %d", i, m.Len(), len(want))
		}
	}

	for k, v := range want {
		got, ok := m.Get(k)
		if !ok || got != v {
			t.Fatalf("Get(%q) = (%v, %v), want (%v, true)", k, got, ok, v)
		}
	}

	visited := 0
	m.Range(func(key string, value interface{}) bool {
		v, ok := want[key]
		if !ok {
			t.Fatalf("Range yielded unexpected key %q", key)
		}
		if value != v {
			t.Fatalf("Range yielded %q = %v, want %v", key, value, v)
		}
		visited++
		return true
	})
	if visited != len(want) {
		t.Fatalf("Range visited %d entries, want %d", visited, len(want))
	}
}

// ExampleHashMap 演示基本用法，可用 `go test -run Example -v` 查看输出。
func ExampleHashMap() {
	m := New()
	m.Set("name", "Alice")
	m.Set("age", 30)

	if v, ok := m.Get("name"); ok {
		fmt.Println("name =", v)
	}

	m.Set("age", 31) // 更新已有 key，不增加 Len
	fmt.Println("len =", m.Len())

	m.Delete("age")
	fmt.Println("len after delete =", m.Len())

	visits := 0
	m.Range(func(key string, value interface{}) bool {
		visits++
		return false // 提前终止
	})
	fmt.Println("range visits =", visits)

	// Output:
	// name = Alice
	// len = 2
	// len after delete = 1
	// range visits = 1
}

func BenchmarkSet(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		m.Set(fmt.Sprintf("key-%d", i), i)
	}
}

func BenchmarkGet(b *testing.B) {
	m := New()
	const n = 10000
	for i := 0; i < n; i++ {
		m.Set(fmt.Sprintf("key-%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(fmt.Sprintf("key-%d", i%n))
	}
}
