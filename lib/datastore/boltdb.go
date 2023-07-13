package datastore

import (
	"encoding/binary"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/encoding"
	"github.com/boltdb/bolt"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
	"sync"
)

// Define byte slices for bucket names
var (
	bucketLogs = []byte("logs")
	bucketConf = []byte("conf")
)

// BoltDbStore is a struct that implements the raft.LogStore interface
// It uses BoltDB as the underlying storage and a read-write mutex for concurrency control
type BoltDbStore struct {
	mux  sync.RWMutex
	conn *bolt.DB
}

// NewLogBoltDbStorage is a function that creates a new BoltDB store
// It takes a configuration map as input and returns a raft.LogStore and an error
func NewLogBoltDbStorage(conf config.Config) (raft.LogStore, error) {
	filename := conf.LogDataStoragePath
	dbOpts := &bolt.Options{
		ReadOnly: false, // Open the database in read-write mode
	}
	conn, err := bolt.Open(filename, 0600, dbOpts) // Open the BoltDB database
	if err != nil {
		return nil, err
	}

	b := &BoltDbStore{
		conn: conn,
	}
	if err := b.init(); err != nil {
		_ = b.conn.Close()
		return nil, err
	}
	return b, nil
}

// init is a method on BoltDbStore that initializes the BoltDB store
// It creates the necessary buckets if they don't exist
func (ds *BoltDbStore) init() error {
	ds.mux.RLock()
	defer ds.mux.RUnlock()
	tx, err := ds.conn.Begin(true)
	if err != nil {
		return err
	}
	defer check(tx.Rollback)

	// Create necessary buckets if they don't exist
	if _, err := tx.CreateBucketIfNotExists(bucketLogs); err != nil {
		return err
	}
	if _, err := tx.CreateBucketIfNotExists(bucketConf); err != nil {
		return err
	}

	return tx.Commit()
}

// FirstIndex is a method on BoltDbStore that returns the first index in the log
func (ds *BoltDbStore) FirstIndex() (uint64, error) {
	tx, err := ds.conn.Begin(false)
	if err != nil {
		return 0, err
	}
	defer check(tx.Rollback)
	var (
		key []byte
		idx uint64
	)
	curs := tx.Bucket(bucketLogs).Cursor()
	key, _ = curs.First() // Retrieve the first key in the bucket
	if key != nil {
		idx = binary.BigEndian.Uint64(key)
	}

	return idx, nil
}

// LastIndex is a method on BoltDbStore that returns the last index in the log
func (ds *BoltDbStore) LastIndex() (uint64, error) {
	tx, err := ds.conn.Begin(false)
	if err != nil {
		return 0, err
	}
	defer check(tx.Rollback)
	var (
		key []byte
		idx uint64
	)
	curs := tx.Bucket(bucketLogs).Cursor() // Get a cursor for the logs bucket
	key, _ = curs.Last()                   // Retrieve the last key in the bucket
	if key != nil {
		idx = binary.BigEndian.Uint64(key) // Convert the key from bytes to uint64
	}

	return idx, nil
}

// GetLog is a method on BoltDbStore that retrieves a log entry by its index
func (ds *BoltDbStore) GetLog(index uint64, log *raft.Log) error {
	ds.mux.RLock()
	defer ds.mux.RUnlock()
	tx, err := ds.conn.Begin(false) // Start a read-only transaction
	if err != nil {
		return err
	}
	defer check(tx.Rollback)

	bucket := tx.Bucket(bucketLogs)         // Get the logs bucket
	val := bucket.Get(uint64ToBytes(index)) // Retrieve the log entry by index

	if val == nil {
		return raft.ErrLogNotFound // Return an error if the log entry is not found
	}
	return encoding.DecodeMessagePack(val, log) // Decode the log entry and assign it to the log variable
}

// StoreLog is a method on BoltDbStore that stores a single log entry
func (ds *BoltDbStore) StoreLog(log *raft.Log) error {
	return ds.StoreLogs([]*raft.Log{log}) // Call the StoreLogs method with a single log entry
}

// StoreLogs is a method on BoltDbStore that stores multiple log entries
func (ds *BoltDbStore) StoreLogs(logs []*raft.Log) error {
	tx, err := ds.conn.Begin(true) // Start a read-write transaction
	if err != nil {
		return err
	}
	defer check(tx.Rollback)
	bucket := tx.Bucket(bucketLogs)
	for _, log := range logs {
		var (
			key []byte
			val []byte
		)
		key = uint64ToBytes(log.Index)             // Convert the index to bytes
		val, err = encoding.EncodeMessagePack(log) // Encode the log entry
		if err != nil {
			break
		}
		err = bucket.Put(key[:], val)
		if err != nil {
			return err
		}
	}

	return tx.Commit() // Commit the transaction
}

// DeleteRange is a method on BoltDbStore that deletes a range of log entries
func (ds *BoltDbStore) DeleteRange(min, max uint64) error {
	tx, err := ds.conn.Begin(true) // Start a read-write transaction
	if err != nil {
		return err
	}
	key := uint64ToBytes(min)

	defer check(tx.Rollback)
	curs := tx.Bucket(bucketLogs).Cursor()
	for k, _ := curs.Seek(key[:]); k != nil; k, _ = curs.Next() {
		// If we reach the max, we are done
		if binary.BigEndian.Uint64(k) > max {
			break
		}
		if err := curs.Delete(); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// uint64ToBytes is a helper function that converts an uint64 to a byte slice
func uint64ToBytes(val uint64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], val)
	return b[:]
}

// uint64ToBytes is a helper function that converts an uint64 to a byte slice
func bytesToUint64(val []byte) uint64 {
	return binary.BigEndian.Uint64(val)
}

func check(f func() error) {
	if err := f(); err != nil {
		// todo create a central logging solution
		zap.L().Log(zap.ErrorLevel, fmt.Sprintf("Received error:%s", err))
	}
}
