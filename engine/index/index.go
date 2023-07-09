package index

import (
	"bytes"
	"github.com/ByteStorage/FlyDB/engine/data"
	"github.com/google/btree"
)

/*
Indexer index interface abstraction layer.
If you want to access other data structures,
you can directly implement this interface
*/
type Indexer interface {
	// Put stores the position information of the key in the index.
	Put(key []byte, pst *data.LogRecordPst) bool

	// Get retrieves the position information of the key from the index.
	Get(key []byte) *data.LogRecordPst

	// Delete deletes the position information of the key from the index.
	Delete(key []byte) bool

	// Size returns the number of entries in the index.
	Size() int

	// Iterator returns an iterator for the index.
	Iterator(reverse bool) Iterator
}

type IndexType = int8

const (
	// Btree Index
	Btree IndexType = iota + 1

	// ART Index
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

// Iterator is a generic index iterator.
type Iterator interface {
	// Rewind resets the iterator to the beginning, i.e., the first entry.
	Rewind()

	// Seek seeks to a target key that is >= or <= the given key, depending on the implementation.
	Seek(key []byte)

	// Next moves to the next key.
	Next()

	// Valid returns whether the iterator is still valid, i.e., if all keys have been traversed.
	Valid() bool

	// Key returns the key at the current iterator position.
	Key() []byte

	// Value returns the value (position information) at the current iterator position.
	Value() *data.LogRecordPst

	// Close closes the iterator and releases any resources.
	Close()
}

func Compare(a, b []byte) int {
	return bytes.Compare(a, b)
}
