package memory

import (
	"errors"
	"sync"
)

type MemTable struct {
	table    map[string][]byte
	size     int64
	mutex    sync.RWMutex
	hasFlush map[string]bool
}

func NewMemTable() *MemTable {
	return &MemTable{
		table:    make(map[string][]byte),
		hasFlush: make(map[string]bool),
	}
}

func (m *MemTable) Put(key string, value []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.table[key] = value
}

func (m *MemTable) Get(key string) ([]byte, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value, ok := m.table[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	return value, nil
}

func (m *MemTable) Flush(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.hasFlush[key] = true
}
