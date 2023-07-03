package index

import (
	"bytes"
	"github.com/ByteStorage/FlyDB/engine/data"
	"github.com/google/btree"
	"sort"
	"sync"
)

/*
	BTree index, which mainly encapsulates Google's btree library
*/

type BTree struct {
	// Source code: Not thread-safe for writing, requires locking; reading does not require locking
	tree *btree.BTree

	lock *sync.RWMutex
}

// NewBTree initializes a new BTree.
func NewBTree() *BTree {
	return &BTree{
		tree: btree.New(32),
		lock: new(sync.RWMutex),
	}
}

func (bt *BTree) Put(key []byte, pst *data.LogRecordPst) bool {
	bt.lock.Lock()
	defer bt.lock.Unlock()

	it := &Item{key: key, pst: pst}
	bt.tree.ReplaceOrInsert(it)
	return true
}

func (bt *BTree) Get(key []byte) *data.LogRecordPst {
	it := &Item{key: key}
	btreeItem := bt.tree.Get(it)
	if btreeItem == nil {
		return nil
	}
	return btreeItem.(*Item).pst
}

func (bt *BTree) Delete(key []byte) bool {
	bt.lock.Lock()
	defer bt.lock.Unlock()

	it := &Item{key: key}
	oldItem := bt.tree.Delete(it)
	if oldItem == nil {
		return false
	}
	return true
}

func (bt *BTree) Size() int {
	return bt.tree.Len()
}

func (bt *BTree) Iterator(reverse bool) Iterator {
	if bt.tree == nil {
		return nil
	}
	bt.lock.RLock()
	defer bt.lock.RUnlock()
	return NewBTreeIterator(bt.tree, reverse)
}

// BTreeIterator represents an iterator for BTree index.
type BtreeIterator struct {
	currIndex int     // Current index position during iteration
	reverse   bool    // Whether it is a reverse iteration
	values    []*Item // Key + position index information
}

func NewBTreeIterator(tree *btree.BTree, reverse bool) *BtreeIterator {
	var idx int
	values := make([]*Item, tree.Len())

	// Store all the data in an array
	saveToValues := func(item btree.Item) bool {
		values[idx] = item.(*Item)
		idx++
		return true
	}
	// Determine whether to traverse in reverse
	if reverse {
		tree.Descend(saveToValues)
	} else {
		tree.Ascend(saveToValues)
	}

	return &BtreeIterator{
		currIndex: 0,
		reverse:   reverse,
		values:    values,
	}

}

func (bi *BtreeIterator) Rewind() {
	bi.currIndex = 0
}

func (bi *BtreeIterator) Seek(key []byte) {
	// Binary search
	if bi.reverse {
		bi.currIndex = sort.Search(len(bi.values), func(i int) bool {
			return bytes.Compare(bi.values[i].key, key) <= 0
		})
	} else {
		bi.currIndex = sort.Search(len(bi.values), func(i int) bool {
			return bytes.Compare(bi.values[i].key, key) >= 0
		})
	}
}

func (bi *BtreeIterator) Next() {
	bi.currIndex += 1
}

func (bi *BtreeIterator) Valid() bool {
	return bi.currIndex < len(bi.values)
}

func (bi *BtreeIterator) Key() []byte {
	return bi.values[bi.currIndex].key
}

func (bi *BtreeIterator) Value() *data.LogRecordPst {
	return bi.values[bi.currIndex].pst
}

func (bi *BtreeIterator) Close() {
	bi.values = nil
}
