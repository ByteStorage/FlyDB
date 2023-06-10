package index

import (
	"bytes"
	"github.com/ByteStorage/flydb/data"
	"github.com/google/btree"
	"sort"
	"sync"
)

/*
BTree 索引，主要封装了 google 的btree库
*/

type BTree struct {
	tree *btree.BTree //源码：多线程写不安全，要加锁;读不需要
	lock *sync.RWMutex
}

// NewBTree 初始化BTree
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

// BTree 索引迭代器
type BtreeIterator struct {
	currIndex int     // 当前遍历的下标位置
	reverse   bool    // 是否是反向遍历
	values    []*Item // key + 位置索引信息
}

func NewBTreeIterator(tree *btree.BTree, reverse bool) *BtreeIterator {
	var idx int
	values := make([]*Item, tree.Len())

	// 将所有的数据存放到数组中
	saveToValues := func(item btree.Item) bool {
		values[idx] = item.(*Item)
		idx++
		return true
	}
	// 判断是否反向遍历
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
	// 二分查找
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
