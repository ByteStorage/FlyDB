package flydb

import (
	"bytes"
	"github.com/ByteStorage/flydb/config"
	"github.com/ByteStorage/flydb/index"
)

// Iterator 迭代器
type Iterator struct {
	indexIter index.Iterator
	db        *DB
	options   config.IteratorOptions
}

// NewIterator 初始化迭代器
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

// 根据Prefix 传进来的 key 判断是与迭代器里面的 key 的前缀相等，不相等往后迭代
func (it *Iterator) skipToNext() {
	prefixLen := len(it.options.Prefix)
	if prefixLen == 0 {
		return
	}

	for ; it.indexIter.Valid(); it.indexIter.Next() {
		key := it.indexIter.Key()
		if prefixLen <= len(key) && bytes.Compare(it.options.Prefix, key[:prefixLen]) == 0 {
			break
		}
	}
}
