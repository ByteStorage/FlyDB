package flydb

import (
	"encoding/binary"
	"github.com/qishenonly/flydb/data"
	"sync"
	"sync/atomic"
)

// 默认非事务id
var nonTransactionSeqNo uint64 = 0

// 标识事务完成
var lgrTransFinaKey = []byte("lgr-fina")

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
		temporaryDataWrites: make(map[string]*data.LogRecord),
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

// Commit 事务提交，将暂存的数据写入数据文件，并更新内存索引
func (wb *WriteBatch) Commit() error {
	wb.lock.Lock()
	defer wb.lock.Unlock()

	if len(wb.temporaryDataWrites) == 0 {
		return nil
	}
	if uint(len(wb.temporaryDataWrites)) > wb.options.MaxBatchNum {
		return ErrExceedMaxBatchNum
	}

	// 获取当前最新的事务序列号
	transSeq := atomic.AddUint64(&wb.db.transSeqNo, 1)

	// 开始写数据到数据文件当中
	// 单条数据写完不会立刻更新索引，需要暂存
	positions := make(map[string]*data.LogRecordPst)
	for _, record := range wb.temporaryDataWrites {
		logRecordPst, err := wb.db.appendLogRecord(&data.LogRecord{
			Key:   encodeLogRecordKeyWithSeq(record.Key, transSeq),
			Value: record.Value,
			Type:  record.Type,
		})
		if err != nil {
			return err
		}
		positions[string(record.Key)] = logRecordPst
	}

	// 写一条标识事务完成的数据
	finishedRecord := &data.LogRecord{
		Key:  encodeLogRecordKeyWithSeq(lgrTransFinaKey, transSeq),
		Type: data.LogRecordTransFinished,
	}
	if _, err := wb.db.appendLogRecord(finishedRecord); err != nil {
		return err
	}

	// 根据配置决定是否持久化
	if wb.options.SyncWrites && wb.db.activeFile != nil {
		if err := wb.db.activeFile.Sync(); err != nil {
			return err
		}
	}

	// 更新内存索引
	for _, record := range wb.temporaryDataWrites {
		pst := positions[string(record.Key)]
		if record.Type == data.LogRecordNormal {
			wb.db.index.Put(record.Key, pst)
		}
		if record.Type == data.LogRecordDeleted {
			wb.db.index.Delete(record.Key)
		}
	}

	// 清空暂存数据
	wb.temporaryDataWrites = make(map[string]*data.LogRecord)

	return nil

}

// encodeLogRecordKeyWithSeq Key+Seq Number 编码
func encodeLogRecordKeyWithSeq(key []byte, seqNo uint64) []byte {
	seq := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(seq[:], seqNo)

	encodeKey := make([]byte, n+len(key))
	copy(encodeKey[:n], seq[:n])
	copy(encodeKey[n:], key)

	return encodeKey
}

// 解析 LogRecord 的 key，获取实际的 key 和事务序列号 seq
func parseLogRecordKeyAndSeq(key []byte) ([]byte, uint64) {
	seqNo, n := binary.Uvarint(key)
	realKey := key[n:]
	return realKey, seqNo
}
