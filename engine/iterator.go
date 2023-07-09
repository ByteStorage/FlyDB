package engine

import (
	"bytes"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine/index"
)

// Iterator iterator
type Iterator struct {
	indexIter index.Iterator
	db        *DB
	options   config.IteratorOptions
}

// NewIterator Initializes the iterator
func (db *DB) NewIterator(opt config.IteratorOptions) *Iterator {
	indexIter := db.index.Iterator(opt.Reverse)
	return &Iterator{
		indexIter: indexIter,
		db:        db,
		options:   opt,
	}
}

func (it *Iterator) Rewind() {
	it.indexIter.Rewind()
	it.skipToNext()
}

func (it *Iterator) Seek(key []byte) {
	it.indexIter.Seek(key)
	it.skipToNext()
}

func (it *Iterator) Next() {
	it.indexIter.Next()
	it.skipToNext()
}

func (it *Iterator) Valid() bool {
	return it.indexIter.Valid()
}

func (it *Iterator) Key() []byte {
	return it.indexIter.Key()
}

func (it *Iterator) Value() ([]byte, error) {
	logRecordPst := it.indexIter.Value()
	it.db.lock.RLock()
	defer it.db.lock.RUnlock()

	return it.db.getValueByPosition(logRecordPst)
}

func (it *Iterator) Close() {
	it.indexIter.Close()
}

// According to the key passed by Prefix,
// the key is equal to the key prefix in the iterator.
// If it is not equal, the next iteration is performed
func (it *Iterator) skipToNext() {
	prefixLen := len(it.options.Prefix)
	if prefixLen == 0 {
		return
	}

	for ; it.indexIter.Valid(); it.indexIter.Next() {
		key := it.indexIter.Key()

		// Check if the key has the desired prefix
		if prefixLen <= len(key) && bytes.Equal(it.options.Prefix, key[:prefixLen]) {
			break
		}
	}
}
