package datastore

import (
	"github.com/ByteStorage/FlyDB/config"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/encoding"
	"github.com/hashicorp/raft"
	"github.com/linxGnu/grocksdb"
	"sync"
)

// RocksDbStore is a struct that implements the raft.LogStore interface
// It uses BoltDB as the underlying storage and a read-write mutex for concurrency control
type RocksDbStore struct {
	mux  sync.RWMutex
	conn *grocksdb.DB
}

// NewLogRocksDbStorage is a function that creates a new RocksDB store
// It takes a configuration map as input and returns a raft.LogStore and an error
func NewLogRocksDbStorage(conf config.Config) (DataStore, error) {
	filename := conf.LogDataStoragePath
	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(grocksdb.NewLRUCache(3 << 30))
	options := grocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)
	options.SetBlockBasedTableFactory(bbto)
	db, err := grocksdb.OpenDb(options, filename)
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
	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()
	it := rds.conn.NewIterator(ro)
	it.SeekToFirst()
	return bytesToUint64(it.Key().Data()), nil
}

// LastIndex is a method on RocksDbStore that returns the last index in the log
func (rds *RocksDbStore) LastIndex() (uint64, error) {
	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()
	it := rds.conn.NewIterator(ro)
	it.SeekToLast()
	return bytesToUint64(it.Key().Data()), nil
}

// GetLog is a method on RocksDbStore that retrieves a log entry by its index
func (rds *RocksDbStore) GetLog(index uint64, log *raft.Log) error {
	ro := grocksdb.NewDefaultReadOptions()
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
	wo := grocksdb.NewDefaultWriteOptions()
	wb := grocksdb.NewWriteBatch()
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
	ro := grocksdb.NewDefaultReadOptions()
	wo := grocksdb.NewDefaultWriteOptions()
	wb := grocksdb.NewWriteBatch()
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

func (rds *RocksDbStore) Set(key []byte, val []byte) error {
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}
	wo := grocksdb.NewDefaultWriteOptions()
	wo.SetSync(true)
	defer wo.Destroy()
	err := rds.conn.Put(wo, key, val)
	if err != nil {
		return err
	}
	return nil
}

func (rds *RocksDbStore) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, _const.ErrKeyIsEmpty
	}
	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()

	val, err := rds.conn.Get(ro, key)
	if err != nil {
		return nil, err
	}
	if val.Size() == 0 {
		return nil, _const.ErrKeyNotFound
	}
	return val.Data(), nil
}

func (rds *RocksDbStore) SetUint64(key []byte, val uint64) error {
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}
	wo := grocksdb.NewDefaultWriteOptions()
	wo.SetSync(true)
	defer wo.Destroy()
	err := rds.conn.Put(wo, key, uint64ToBytes(val))
	if err != nil {
		return err
	}
	return nil
}

func (rds *RocksDbStore) GetUint64(key []byte) (uint64, error) {
	if len(key) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}
	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()

	val, err := rds.conn.Get(ro, key)
	if err != nil {
		return 0, err
	}
	return bytesToUint64(val.Data()), nil
}
