package storage

import (
	"sync"
	"time"
)

type Value struct {
	Data   []byte
	Expiry int64
}

type MemTable struct {
	mu   sync.RWMutex
	data map[string]Value
}

func NewMemTable() *MemTable {
	return &MemTable{
		data: make(map[string]Value),
	}
}

func (m *MemTable) Get(key string) ([]byte, bool) {
	m.mu.RLock()
	v, ok := m.data[key]
	m.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if v.Expiry > 0 && time.Now().UnixNano() > v.Expiry {
		m.mu.Lock()
		delete(m.data, key)
		m.mu.Unlock()
		return nil, false
	}

	return v.Data, true
}

func (m *MemTable) Put(key string, value []byte, ttl time.Duration) {
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	m.mu.Lock()
	m.data[key] = Value{
		Data:   value,
		Expiry: exp,
	}
	m.mu.Unlock()
}

func (m *MemTable) Delete(key string) {
	m.mu.Lock()
	delete(m.data, key)
	m.mu.Unlock()
}
