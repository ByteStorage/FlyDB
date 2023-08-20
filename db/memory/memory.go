package memory

import (
	"errors"
	"github.com/ByteStorage/FlyDB/lib/bloom"
	"sync"
)

// MemTable is an in-memory table
type MemTable struct {
	table map[string][]byte // key -> value
	mutex sync.RWMutex      // protect table
	bloom *bloom.Filter     // bloom filter
}

// NewMemTable create a new MemTable
func NewMemTable() *MemTable {
	return &MemTable{
		table: make(map[string][]byte),
		// Initialize with no keys and 10 bits per key
		bloom: bloom.NewBloomFilter(1000, 0.01),
	}
}

// Put a key-value pair into the table
func (m *MemTable) Put(key string, value []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.table[key] = value
	// Add the key to the bloom filter
	m.bloom.Add([]byte(key))
}

// Get a value from the table
func (m *MemTable) Get(key string) ([]byte, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Immediate return if the key is not in the bloom filter
	if !m.bloom.MayContainItem([]byte(key)) {
		return nil, errors.New("key not found")
	}

	value, ok := m.table[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	return value, nil
}

// Delete a key from the table
// Note: Bloom filters don't support deletion without affecting accuracy
// so we don't remove the key from the bloom filter.
func (m *MemTable) Delete(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.table, key)
}
