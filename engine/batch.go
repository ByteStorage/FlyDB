package engine

import (
	"encoding/binary"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine/data"
	"github.com/ByteStorage/FlyDB/lib/const"
	"sync"
	"sync/atomic"
)

// Default non-transaction id
var nonTransactionSeqNo uint64 = 0

// Identify transaction completion
var lgrTransFinaKey = []byte("lgr-fina")

// WriteBatch Writes data in atomic batches to ensure atomicity
type WriteBatch struct {
	options             config.WriteBatchOptions
	lock                *sync.Mutex
	db                  *DB
	temporaryDataWrites map[string]*data.LogRecord // Stores the data written by the user
}

// NewWriteBatch Init WriteBatch
func (db *DB) NewWriteBatch(opt config.WriteBatchOptions) *WriteBatch {
	return &WriteBatch{
		options:             opt,
		lock:                new(sync.Mutex),
		db:                  db,
		temporaryDataWrites: make(map[string]*data.LogRecord),
	}
}

// Put Data batch write
func (wb *WriteBatch) Put(key []byte, value []byte) error {
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}
	wb.lock.Lock()
	defer wb.lock.Unlock()

	// Temporarily store the LogRecord
	logRecord := &data.LogRecord{
		Key:   key,
		Value: value,
	}
	wb.temporaryDataWrites[string(key)] = logRecord
	return nil
}

// Delete Batch deletion of data
func (wb *WriteBatch) Delete(key []byte) error {
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}
	wb.lock.Lock()
	defer wb.lock.Unlock()

	// If the data does not exist, return it directly
	logRecordPst := wb.db.index.Get(key)
	if logRecordPst == nil {
		if wb.temporaryDataWrites[string(key)] != nil {
			delete(wb.temporaryDataWrites, string(key))
		}
		return nil
	}

	// Temporarily store the LogRecord
	logRecord := &data.LogRecord{
		Key:  key,
		Type: data.LogRecordDeleted,
	}
	wb.temporaryDataWrites[string(key)] = logRecord
	return nil
}

// Commit The transaction commits, writes the transient data to the data file, and updates the in-memory index
func (wb *WriteBatch) Commit() error {
	wb.lock.Lock()
	defer wb.lock.Unlock()

	if len(wb.temporaryDataWrites) == 0 {
		return nil
	}
	if uint(len(wb.temporaryDataWrites)) > wb.options.MaxBatchNum {
		return _const.ErrExceedMaxBatchNum
	}

	// Gets the current, most recent transaction sequence number
	transSeq := atomic.AddUint64(&wb.db.transSeqNo, 1)

	// Start writing data to the data file
	// The index is not updated immediately after a single piece of data is written.
	// It needs to be stored temporarily
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

	// Write a piece of data that identifies the completion of the transaction
	finishedRecord := &data.LogRecord{
		Key:  encodeLogRecordKeyWithSeq(lgrTransFinaKey, transSeq),
		Type: data.LogRecordTransFinished,
	}
	if _, err := wb.db.appendLogRecord(finishedRecord); err != nil {
		return err
	}

	// Decide whether to persist based on the configuration
	if wb.options.SyncWrites && wb.db.activeFile != nil {
		if err := wb.db.activeFile.Sync(); err != nil {
			return err
		}
	}

	// Update memory index
	for _, record := range wb.temporaryDataWrites {
		pst := positions[string(record.Key)]
		if record.Type == data.LogRecordNormal {
			wb.db.index.Put(record.Key, pst)
		}
		if record.Type == data.LogRecordDeleted {
			wb.db.index.Delete(record.Key)
		}
	}

	// Clear the temporary data
	wb.temporaryDataWrites = make(map[string]*data.LogRecord)

	return nil

}

// encodeLogRecordKeyWithSeq Key+Seq Number coding
func encodeLogRecordKeyWithSeq(key []byte, seqNo uint64) []byte {
	seq := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(seq[:], seqNo)

	encodeKey := make([]byte, n+len(key))
	copy(encodeKey[:n], seq[:n])
	copy(encodeKey[n:], key)

	return encodeKey
}

// Parse the LogRecord key to get the actual key and transaction sequence number seq
func parseLogRecordKeyAndSeq(key []byte) ([]byte, uint64) {
	seqNo, n := binary.Uvarint(key)
	realKey := key[n:]
	return realKey, seqNo
}
