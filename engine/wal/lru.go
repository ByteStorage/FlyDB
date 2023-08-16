package wal

import (
	"container/list"
	"sync"
)

type Cache struct {
	capacity int
	ll       *list.List
	mu       sync.Mutex // Mutex to ensure concurrent safety for linked list operations.
	cache    sync.Map   // Using sync.Map to safely handle concurrent cache operations.
}

type entry struct {
	key   string // Using string as the key type internally to avoid issues with []byte being unhashable.
	value []byte // Using []byte as the value type.
}

// NewCache initializes a new LRU Cache with the given capacity.
func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		ll:       list.New(),
	}
}

// Get retrieves a value from the cache using a []byte key.
func (c *Cache) Get(key []byte) (value []byte, ok bool) {
	return c.get(string(key))
}

// Put adds a key-value pair to the cache, using a []byte key and value.
func (c *Cache) Put(key, value []byte) {
	c.put(string(key), value)
}

// Internal get function using string as the key.
func (c *Cache) get(key string) (value []byte, ok bool) {
	v, ok := c.cache.Load(key)
	if !ok {
		return nil, false
	}

	// Move the accessed item to the front of the linked list.
	c.mu.Lock()
	ele := v.(*list.Element)
	c.ll.MoveToFront(ele)
	c.mu.Unlock()

	return ele.Value.(*entry).value, true
}

// Internal put function using string as the key.
func (c *Cache) put(key string, value []byte) {
	// Check if the key already exists in the cache.
	if v, ok := c.cache.Load(key); ok {
		c.mu.Lock()
		ele := v.(*list.Element)
		ele.Value.(*entry).value = value
		c.ll.MoveToFront(ele)
		c.mu.Unlock()
		return
	}

	e := &entry{key, value}
	c.mu.Lock()
	ele := c.ll.PushFront(e)
	if c.ll.Len() > c.capacity {
		c.removeOldest()
	}
	c.mu.Unlock()
	c.cache.Store(key, ele)
}

// removeOldest removes the least recently used item from the cache.
func (c *Cache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		c.cache.Delete(kv.key)
	}
}
