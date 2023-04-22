package flydb

import (
	"github.com/qishenonly/flydb/data"
	"sync"
)

// WriteBatch 原子批量写数据，保证操作原子性
type WriteBatch struct {
	options             WriteBatchOptions
	lock                *sync.Mutex
	db                  *DB
	temporaryDataWrites map[string]*data.LogRecord // 暂存用户写入的数据
}

// NewWriteBatch 初始化 WriteBatch
func (db *DB) NewWriteBatch(opt WriteBatchOptions) *WriteBatch {
	return &WriteBatch{
		options:             opt,
		lock:                new(sync.Mutex),
		db:                  db,
		temporaryDataWrites: map[string]*data.LogRecord{},
	}
}

// Put 数据批量写入
func (wb *WriteBatch) Put(key []byte, value []byte) error {
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}
	wb.lock.Lock()
	defer wb.lock.Unlock()

	// 暂存 LogRecord
	logRecord := &data.LogRecord{
		Key:   key,
		Value: value,
	}
	wb.temporaryDataWrites[string(key)] = logRecord
	return nil
}

// Delete 数据批量删除
func (wb *WriteBatch) Delete(key []byte) error {
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}
	wb.lock.Lock()
	defer wb.lock.Unlock()

	// 数据不存在则直接返回
	logRecordPst := wb.db.index.Get(key)
	if logRecordPst == nil {
		if wb.temporaryDataWrites[string(key)] != nil {
			delete(wb.temporaryDataWrites, string(key))
		}
		return nil
	}

	// 暂存 LogRecord
	logRecord := &data.LogRecord{
		Key:  key,
		Type: data.LogRecordDeleted,
	}
	wb.temporaryDataWrites[string(key)] = logRecord
	return nil
}


