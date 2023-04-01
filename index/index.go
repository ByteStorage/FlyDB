package index

import (
	"bytes"
	"flydb/data"
	"github.com/google/btree"
)

/*
Indexer 索引接口抽象层，后续如果想要接入其他的数据结构，则直接实现这个接口
*/
type Indexer interface {
	// Put 向索引中存储 key 对应的数据位置信息
	Put(key []byte, pst *data.LogRecordPst) bool

	// Get 根据 key 取出对应的索引位置信息
	Get(key []byte) *data.LogRecordPst

	// Delete 根据 key 删除对应的索引位置信息
	Delete(key []byte) bool
}

type Item struct {
	key []byte
	pst *data.LogRecordPst
}

func (i *Item) Less(bi btree.Item) bool {
	return bytes.Compare(i.key, bi.(*Item).key) == -1
}
