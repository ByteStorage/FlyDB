package index

import (
	art "github.com/plar/go-adaptive-radix-tree"
	"github.com/qishenonly/flydb/data"
	"sync"
)

// Adaptive Radix Tree Index
// https://github.com/plar/go-adaptive-radix-tree
type AdaptiveRadixTree struct {
	tree art.Tree
	lock *sync.RWMutex
}

// NewART Initializes the adaptive radix tree index
func NewART() *AdaptiveRadixTree {
	return &AdaptiveRadixTree{
		tree: art.New(),
		lock: new(sync.RWMutex),
	}
}

func (artree *AdaptiveRadixTree) Put(key []byte, pst *data.LogRecordPst) bool {
	artree.lock.Lock()
	defer artree.lock.Unlock()
	artree.tree.Insert(key, pst)
	return true
}

func (artree *AdaptiveRadixTree) Get(key []byte) *data.LogRecordPst {
	artree.lock.RLock()
	defer artree.lock.RUnlock()
	value, found := artree.tree.Search(key)
	if !found {
		return nil
	}
	return value.(*data.LogRecordPst)
}

func (artree *AdaptiveRadixTree) Delete(key []byte) bool {
	artree.lock.Lock()
	defer artree.lock.Unlock()
	_, deleted := artree.tree.Delete(key)
	return deleted
}

func (artree *AdaptiveRadixTree) Size() int {
	artree.lock.RLock()
	defer artree.lock.RUnlock()
	size := artree.tree.Size()
	return size
}

func (artree *AdaptiveRadixTree) Iterator(reverse bool) Iterator {

}
