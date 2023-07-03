package index

import (
	"bytes"
	"github.com/ByteStorage/FlyDB/data"
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

	// Size 索引中的数据量
	Size() int

	// Iterator 索引迭代器
	Iterator(reverse bool) Iterator
}

type IndexType = int8

const (
	// Btree 索引
	Btree IndexType = iota + 1

	// ART 自适应基数数索引
	ART
)

func NewIndexer(typeIndex IndexType, dirPath string) Indexer {
	switch typeIndex {
	case Btree:
		return NewBTree()
	case ART:
		return NewART()
	default:
		panic("unsupported index type")
	}
}

type Item struct {
	key []byte
	pst *data.LogRecordPst
}

func (i *Item) Less(bi btree.Item) bool {
	return bytes.Compare(i.key, bi.(*Item).key) == -1
}

// Iterator 通用索引迭代器
type Iterator interface {
	// Rewind 重新回到迭代器的起点，即第一个数据
	Rewind()

	// Seek 根据传入的 key 查找到一个 >= 或 <= 的目标 key，从这个目标 key 开始遍历
	Seek(key []byte)

	// Next 跳转到下一个 key
	Next()

	// Valid 是否有效，即是否已经遍历完了所有 key，用于退出遍历 ==> true->是  false-->否
	Valid() bool

	// Key 当前遍历位置的 key 数据
	Key() []byte

	// Value 当前遍历位置的 value 数据
	Value() *data.LogRecordPst

	// Close 关闭迭代器，释放相应资源
	Close()
}
