package memory

import (
	"errors"
	"sync"
)

type MemTable struct {
	table map[string][]byte
	mutex sync.RWMutex
}

func NewMemTable() *MemTable {
	return &MemTable{
		table: make(map[string][]byte),
	}
}

func (m *MemTable) Put(key []byte, value []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.table[string(key)] = value
}

func (m *MemTable) Get(key []byte) ([]byte, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value, ok := m.table[string(key)]
	if !ok {
		return nil, errors.New("key not found")
	}

	return value, nil
}
