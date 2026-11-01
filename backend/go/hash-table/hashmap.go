package hashmap

const loadFactor = 0.75

type entry struct {
	key   string
	value interface{}
	next  *entry
}

type HashMap struct {
	buckets []*entry
	size    int
}

func New() *HashMap {
	return &HashMap{buckets: make([]*entry, 8)}
}

func hashString(s string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

func (m *HashMap) index(key string) int {
	return int(hashString(key) & uint32(len(m.buckets)-1))
}

func (m *HashMap) Set(key string, value interface{}) {
	idx := m.index(key)
	for e := m.buckets[idx]; e != nil; e = e.next {
		if e.key == key {
			e.value = value
			return
		}
	}

	m.buckets[idx] = &entry{
		key: key, value: value, next: m.buckets[idx],
	}
	m.size++

	if float64(m.size)/float64(len(m.buckets)) > loadFactor {
		m.resize()
	}
}

func (m *HashMap) Get(key string) (interface{}, bool) {
	for e := m.buckets[m.index(key)]; e != nil; e = e.next {
		if e.key == key {
			return e.value, true
		}
	}
	return nil, false
}

func (m *HashMap) Delete(key string) {
	idx := m.index(key)
	var prev *entry

	for e := m.buckets[idx]; e != nil; e = e.next {
		if e.key == key {
			if prev == nil {
				m.buckets[idx] = e.next
			} else {
				prev.next = e.next
			}
			m.size--
			return
		}
		prev = e
	}
}

func (m *HashMap) Len() int {
	return m.size
}

func (m *HashMap) resize() {
	old := m.buckets
	m.buckets = make([]*entry, len(old)*2)

	for _, head := range old {
		for e := head; e != nil; {
			next := e.next
			idx := m.index(e.key)
			e.next = m.buckets[idx]
			m.buckets[idx] = e
			e = next
		}
	}
}

func (m *HashMap) Range(fn func(key string, value interface{}) bool) {
	for _, head := range m.buckets {
		for e := head; e != nil; e = e.next {
			if !fn(e.key, e.value) {
				return
			}
		}
	}
}
