package index

import (
	"github.com/ByteStorage/FlyDB/db/data"
	"github.com/ByteStorage/FlyDB/lib/bloom"
	art "github.com/plar/go-adaptive-radix-tree"
	"sync"
)

// AdaptiveRadixTreeWithBloom Adaptive Radix Tree Index
// The following link is the ART library written by go.
// If you need to know more about it, please go to the corresponding warehouse.
// https://github.com/plar/go-adaptive-radix-tree
type AdaptiveRadixTreeWithBloom struct {
	tree   art.Tree
	lock   *sync.RWMutex
	filter *bloom.Filter
}

// NewARTWithBloom Initializes the adaptive radix tree index
func NewARTWithBloom() *AdaptiveRadixTreeWithBloom {
	return &AdaptiveRadixTreeWithBloom{
		tree:   art.New(),
		lock:   new(sync.RWMutex),
		filter: bloom.NewBloomFilter(1000, 0.01),
	}
}

func (artree *AdaptiveRadixTreeWithBloom) Put(key []byte, pst *data.LogRecordPst) bool {
	artree.lock.Lock()
	defer artree.lock.Unlock()
	artree.tree.Insert(key, pst)
	artree.filter.Add(key)
	return true
}

func (artree *AdaptiveRadixTreeWithBloom) Get(key []byte) *data.LogRecordPst {
	if !artree.filter.MayContainItem(key) {
		return nil
	}
	artree.lock.RLock()
	defer artree.lock.RUnlock()
	value, found := artree.tree.Search(key)
	if !found {
		return nil
	}
	return value.(*data.LogRecordPst)
}

func (artree *AdaptiveRadixTreeWithBloom) Delete(key []byte) bool {
	if !artree.filter.MayContainItem(key) {
		return false
	}
	artree.lock.Lock()
	defer artree.lock.Unlock()
	_, deleted := artree.tree.Delete(key)
	return deleted
}

func (artree *AdaptiveRadixTreeWithBloom) Size() int {
	artree.lock.RLock()
	defer artree.lock.RUnlock()
	size := artree.tree.Size()
	return size
}

func (artree *AdaptiveRadixTreeWithBloom) Iterator(reverse bool) Iterator {
	artree.lock.RLock()
	defer artree.lock.RUnlock()
	return NewARTreeIterator(artree.tree, reverse)
}
