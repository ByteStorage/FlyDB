package index

import (
	"flydb/data"
	"github.com/google/btree"
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
