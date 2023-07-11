package store

import "github.com/hashicorp/raft"

// stableLog implements raft.StableStore interface, we can use it to store raft stable logs
// how to store raft stable logs? we can use FlyDB/RocksDB/LevelDB/BoltDB to store raft stable logs
// maybe we can use FlyDB
type stableLog struct {
	// implement me
}

// newStableLog returns a new stableLog, we can use it to store raft stable logs
func newStableLog() raft.StableStore {
	return &stableLog{}
}

func (s stableLog) Set(key []byte, val []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s stableLog) Get(key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s stableLog) SetUint64(key []byte, val uint64) error {
	//TODO implement me
	panic("implement me")
}

func (s stableLog) GetUint64(key []byte) (uint64, error) {
	//TODO implement me
	panic("implement me")
}
