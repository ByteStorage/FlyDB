package memory

import (
	"errors"
	"sync"
)

// MemTable is a in-memory table
type MemTable struct {
	table map[string][]byte // key -> value
	mutex sync.RWMutex      // protect table
}

// NewMemTable create a new MemTable
func NewMemTable() *MemTable {
	return &MemTable{
		table: make(map[string][]byte),
	}
}

// Put a key-value pair into the table
func (m *MemTable) Put(key string, value []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.table[key] = value
}

// Get a value from the table
func (m *MemTable) Get(key string) ([]byte, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value, ok := m.table[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	return value, nil
}

// Delete a key from the table
func (m *MemTable) Delete(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.table, key)
}
