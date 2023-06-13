package index

import (
	"bytes"
	"github.com/ByteStorage/FlyDB/data"
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

func NewARTreeIterator(tree art.Tree, reverse bool) *BtreeIterator {
	var idx int
	if reverse {
		idx = tree.Size() - 1
	}
	values := make([]*Item, tree.Size())

	// Store all the data in an array
	saveToValues := func(node art.Node) bool {
		item := &Item{
			key: node.Key(),
			pst: node.Value().(*data.LogRecordPst),
		}
		values[idx] = item
		if reverse {
			idx--
		} else {
			idx++
		}
		return true
	}
	tree.ForEach(saveToValues)

	return &BtreeIterator{
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
