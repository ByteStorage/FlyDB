package datastore

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/encoding"
	"github.com/hashicorp/raft"
	"math"
	"sync"
)

// FlyDbStore is a struct that implements the raft.LogStore interface
// It uses FlyDB as the underlying storage and a read-write mutex for concurrency control
type FlyDbStore struct {
	mux        sync.RWMutex
	conn       *engine.DB
	stableConn *engine.DB
	firstIndex uint64
	lastIndex  uint64
}

// NewLogFlyDbStorage is a function that creates a new FlyDB store
// It takes a configuration map as input and returns a raft.LogStore and an error
func NewLogFlyDbStorage(conf config.Config) (DataStore, error) {
	opts := config.DefaultOptions
	opts.DirPath = conf.LogDataStoragePath
	opts.DataFileSize = conf.LogDataStorageSize
	conn, err := engine.NewDB(opts)
	if err != nil {
		return nil, err
	}
	b := &FlyDbStore{
		conn: conn,
	}
	return b, nil
}

// FirstIndex is a method on FlyDbStore that returns the first index in the log
func (fds *FlyDbStore) FirstIndex() (uint64, error) {
	return fds.firstIndex, nil
}

// LastIndex is a method on FlyDbStore that returns the last index in the log
func (fds *FlyDbStore) LastIndex() (uint64, error) {
	return fds.lastIndex, nil
}

// GetLog is a method on FlyDbStore that retrieves a log entry by its index
func (fds *FlyDbStore) GetLog(index uint64, log *raft.Log) error {
	val, err := fds.conn.Get(uint64ToBytes(index)) // Retrieve the log entry by index
	if err != nil {
		return err
	}
	if val == nil {
		return raft.ErrLogNotFound // Return an error if the log entry is not found
	}
	return encoding.DecodeMessagePack(val, log) // Decode the log entry and assign it to the log variable
}

// StoreLog is a method on FlyDbStore that stores a single log entry
func (fds *FlyDbStore) StoreLog(log *raft.Log) error {
	return fds.StoreLogs([]*raft.Log{log}) // Call the StoreLogs method with a single log entry
}

// StoreLogs is a method on FlyDbStore that stores multiple log entries
// since FlyDB currently does not support transactions,
// please be aware that in case of errors, already written data will persist.
func (fds *FlyDbStore) StoreLogs(logs []*raft.Log) error {

	for _, log := range logs {
		var (
			key []byte
			val []byte
		)

		key = uint64ToBytes(log.Index)              // Convert the index to bytes
		val, err := encoding.EncodeMessagePack(log) // Encode the log entry
		if err != nil {
			break
		}
		err = fds.conn.Put(key[:], val)
		if err != nil {
			return err
		}
		if log.Index > fds.lastIndex {
			fds.lastIndex = log.Index
		} else if log.Index < fds.firstIndex || fds.firstIndex == 0 {
			fds.firstIndex = log.Index
		}
	}
	return nil
}

// DeleteRange is a method on FlyDbStore that deletes a range of log entries
func (fds *FlyDbStore) DeleteRange(min, max uint64) error {
	for i := min; i <= max; i++ {
		_ = fds.conn.Delete(uint64ToBytes(i))
	}
	if fds.firstIndex >= min && fds.firstIndex <= max {
		fds.firstIndex = fds.min()
	}
	if fds.lastIndex >= min && fds.lastIndex <= max {
		fds.lastIndex = fds.max()
	}
	return nil
}

// Set is used to store key/value pair
func (fds *FlyDbStore) Set(key []byte, val []byte) error {
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}
	err := fds.conn.Put(key, val)
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves the value associated with `key`; returns []byte if key exists
func (fds *FlyDbStore) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, _const.ErrKeyIsEmpty
	}
	val, err := fds.conn.Get(key)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// SetUint64 is used to set Uint64 value for `key`
func (fds *FlyDbStore) SetUint64(key []byte, val uint64) error {
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}
	err := fds.conn.Put(key, uint64ToBytes(val))
	if err != nil {
		return err
	}
	return nil
}

// GetUint64 is used to retrieve the Uint64 value of `key`
func (fds *FlyDbStore) GetUint64(key []byte) (uint64, error) {
	if len(key) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}
	val, err := fds.conn.Get(key)
	if err != nil {
		return 0, err
	}
	return bytesToUint64(val), nil
}

// min is a helper method on FlyDbStore that returns the smallest index in the log
func (fds *FlyDbStore) min() uint64 {
	var min uint64
	min = math.MaxUint64
	keys := fds.conn.GetListKeys()
	for _, k := range keys {
		if bytesToUint64(k) < min {
			min = bytesToUint64(k)
		}
	}
	return min
}

// max is a helper method on FlyDbStore that returns the largest index in the log
func (fds *FlyDbStore) max() uint64 {
	var max uint64
	keys := fds.conn.GetListKeys()
	for _, k := range keys {
		if bytesToUint64(k) > max {
			max = bytesToUint64(k)
		}
	}
	return max
}
