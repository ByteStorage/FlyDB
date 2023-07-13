package datastore

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/encoding"
	"github.com/hashicorp/raft"
	"github.com/tecbot/gorocksdb"
	"sync"
)

// RocksDbStore is a struct that implements the raft.LogStore interface
// It uses BoltDB as the underlying storage and a read-write mutex for concurrency control
type RocksDbStore struct {
	mux  sync.RWMutex
	conn *gorocksdb.DB
}

// NewLogRocksDbStorage is a function that creates a new RocksDB store
// It takes a configuration map as input and returns a raft.LogStore and an error
func NewLogRocksDbStorage(conf config.Config) (raft.LogStore, error) {
	filename := conf.LogDataStoragePath
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(options, filename)
	if err != nil {
		return nil, err
	}
	b := &RocksDbStore{
		conn: db,
	}

	return b, nil
}

// FirstIndex is a method on RocksDbStore that returns the first index in the log
func (rds *RocksDbStore) FirstIndex() (uint64, error) {
	ro := gorocksdb.NewDefaultReadOptions()
	defer ro.Destroy()
	it := rds.conn.NewIterator(ro)
	it.SeekToFirst()
	return bytesToUint64(it.Key().Data()), nil
}

// LastIndex is a method on RocksDbStore that returns the last index in the log
func (rds *RocksDbStore) LastIndex() (uint64, error) {
	ro := gorocksdb.NewDefaultReadOptions()
	defer ro.Destroy()
	it := rds.conn.NewIterator(ro)
	it.SeekToLast()
	return bytesToUint64(it.Key().Data()), nil
}

// GetLog is a method on RocksDbStore that retrieves a log entry by its index
func (rds *RocksDbStore) GetLog(index uint64, log *raft.Log) error {
	ro := gorocksdb.NewDefaultReadOptions()
	defer ro.Destroy()

	val, err := rds.conn.Get(ro, uint64ToBytes(index))
	if err != nil {
		return err
	}
	err = encoding.DecodeMessagePack(val.Data(), log)
	if err != nil {
		return err
	}
	return nil
}

// StoreLog is a method on RocksDbStore that stores a single log entry
func (rds *RocksDbStore) StoreLog(log *raft.Log) error {
	return rds.StoreLogs([]*raft.Log{log}) // Call the StoreLogs method with a single log entry
}

// StoreLogs is a method on RocksDbStore that stores multiple log entries
func (rds *RocksDbStore) StoreLogs(logs []*raft.Log) error {
	wo := gorocksdb.NewDefaultWriteOptions()
	wb := gorocksdb.NewWriteBatch()
	wo.SetSync(true)
	defer func() {
		wo.Destroy()
		wb.Destroy()
	}()
	for _, log := range logs {
		var (
			key []byte
			val []byte
		)
		key = uint64ToBytes(log.Index)              // Convert the index to bytes
		val, err := encoding.EncodeMessagePack(log) // Encode the log entry
		if err != nil {
			return err
		}
		wb.Put(key, val)

	}

	return rds.conn.Write(wo, wb) // Commit the transaction
}

// DeleteRange is a method on RocksDbStore that deletes a range of log entries
func (rds *RocksDbStore) DeleteRange(min, max uint64) error {
	ro := gorocksdb.NewDefaultReadOptions()
	wo := gorocksdb.NewDefaultWriteOptions()
	wb := gorocksdb.NewWriteBatch()
	defer func() {
		ro.Destroy()
		wo.Destroy()
		wb.Destroy()
	}()
	wb.DeleteRange(uint64ToBytes(min), uint64ToBytes(max))
	err := rds.conn.Write(wo, wb)
	if err != nil {
		return err
	}
	return nil
}
