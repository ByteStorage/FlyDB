package datastore

import (
	"errors"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/hashicorp/raft"
	"math"
	"sync"
)

// InMemStore is a struct that implements the raft.LogStore interface
// It uses a map as the underlying storage and a read-write mutex for concurrency control
type InMemStore struct {
	mux        sync.RWMutex
	firstIndex uint64
	lastIndex  uint64
	logs       map[uint64]*raft.Log
}

// NewLogInMemStorage is a function that creates a new in-memory store
// It returns a raft.LogStore and an error
func NewLogInMemStorage(conf config.Config) (raft.LogStore, error) {
	a := &InMemStore{
		logs: make(map[uint64]*raft.Log),
	}
	return a, nil
}

// FirstIndex is a method on InMemStore that returns the first index in the log
func (ds *InMemStore) FirstIndex() (uint64, error) {
	ds.mux.RLock()
	defer ds.mux.RUnlock()
	return ds.firstIndex, nil
}

// LastIndex is a method on InMemStore that returns the last index in the log
func (ds *InMemStore) LastIndex() (uint64, error) {
	ds.mux.RLock()
	defer ds.mux.RUnlock()
	return ds.lastIndex, nil
}

// GetLog is a method on InMemStore that retrieves a log entry by its index
func (ds *InMemStore) GetLog(index uint64, log *raft.Log) error {
	ds.mux.RLock()
	defer ds.mux.RUnlock()
	l, ok := ds.logs[index]
	if !ok {
		return errors.New("InMem Log Store: log with index not found")
	}
	*log = *l
	return nil
}

// StoreLog is a method on InMemStore that stores a single log entry
func (ds *InMemStore) StoreLog(log *raft.Log) error {
	return ds.StoreLogs([]*raft.Log{log})
}

// StoreLogs is a method on InMemStore that stores multiple log entries
func (ds *InMemStore) StoreLogs(logs []*raft.Log) error {
	ds.mux.Lock()
	defer ds.mux.Unlock()
	for _, l := range logs {
		ds.logs[l.Index] = l
		if ds.firstIndex == 0 {
			ds.firstIndex = l.Index
		}
		if l.Index > ds.lastIndex {
			ds.lastIndex = l.Index
		}
		if l.Index < ds.firstIndex {
			ds.firstIndex = l.Index
		}
	}
	return nil
}

// DeleteRange is a method on InMemStore that deletes a range of log entries
func (ds *InMemStore) DeleteRange(min, max uint64) error {
	ds.mux.Lock()
	defer ds.mux.Unlock()
	for j := min; j <= max; j++ {
		delete(ds.logs, j)
	}
	if ds.firstIndex >= min && ds.firstIndex <= max {
		ds.firstIndex = ds.min()
	}
	if ds.lastIndex >= min && ds.lastIndex <= max {
		ds.lastIndex = ds.max()
	}

	return nil
}

// min is a helper method on InMemStore that returns the smallest index in the log
func (ds *InMemStore) min() uint64 {
	var min uint64
	min = math.MaxUint64
	for k := range ds.logs {
		if k < min {
			min = k
		}
	}
	return min
}

// max is a helper method on InMemStore that returns the largest index in the log
func (ds *InMemStore) max() uint64 {
	var max uint64

	for k := range ds.logs {
		if k > max {
			max = k
		}
	}
	return max
}
