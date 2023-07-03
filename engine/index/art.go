package index

import (
	"bytes"
	"github.com/ByteStorage/FlyDB/engine/data"
	art "github.com/plar/go-adaptive-radix-tree"
	"sort"
	"sync"
)

// Adaptive Radix Tree Index
// The following link is the ART library written by go.
// If you need to know more about it, please go to the corresponding warehouse.
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
	artree.lock.RLock()
	defer artree.lock.RUnlock()
	return NewARTreeIterator(artree.tree, reverse)
}

// ART Index iterator
type ARTreeIterator struct {
	currIndex int     // The subscript position of the current traversal
	reverse   bool    // Whether it is reverse traversal
	values    []*Item // Key + Location index information
}

func NewARTreeIterator(tree art.Tree, reverse bool) *ARTreeIterator {
	// Estimate the expected slice capacity based on tree size
	expectedSize := tree.Size()

	// Use a mutex for concurrent access to values array
	var mutex sync.Mutex

	// Initialize with empty slice and expected capacity
	values := make([]*Item, 0, expectedSize)

	// Store all the data in an array
	saveToValues := func(node art.Node) bool {
		item := &Item{
			key: node.Key(),
			pst: node.Value().(*data.LogRecordPst),
		}
		mutex.Lock()

		// Append item to values slice
		values = append(values, item)

		mutex.Unlock()

		return true
	}
	tree.ForEach(saveToValues)

	// Reverse the values slice if reverse is true
	if reverse {
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
	}

	return &ARTreeIterator{
		currIndex: 0,
		reverse:   reverse,
		values:    values,
	}
}

func (artree *ARTreeIterator) Rewind() {
	artree.currIndex = 0
}

func (artree *ARTreeIterator) Seek(key []byte) {
	// Binary search algorithm
	if artree.reverse {
		artree.currIndex = sort.Search(len(artree.values), func(i int) bool {
			return bytes.Compare(artree.values[i].key, key) <= 0
		})
	} else {
		artree.currIndex = sort.Search(len(artree.values), func(i int) bool {
			return bytes.Compare(artree.values[i].key, key) >= 0
		})
	}
}

func (artree *ARTreeIterator) Next() {
	artree.currIndex += 1
}

func (artree *ARTreeIterator) Valid() bool {
	return artree.currIndex < len(artree.values)
}

func (artree *ARTreeIterator) Key() []byte {
	return artree.values[artree.currIndex].key
}

func (artree *ARTreeIterator) Value() *data.LogRecordPst {
	return artree.values[artree.currIndex].pst
}

func (artree *ARTreeIterator) Close() {
	artree.values = nil
}
