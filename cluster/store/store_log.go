package store

import "github.com/hashicorp/raft"

// raftLog implements raft.LogStore interface, we can use it to store raft logs
// how to store raft logs? we can use FlyDB/RocksDB/LevelDB/BoltDB to store raft logs
// maybe we can use FlyDB
type raftLog struct {
	// implement me
}

// newRaftLog returns a new raftLog
func newRaftLog() raft.LogStore {
	return &raftLog{}
}

func (r *raftLog) FirstIndex() (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *raftLog) LastIndex() (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *raftLog) GetLog(index uint64, log *raft.Log) error {
	//TODO implement me
	panic("implement me")
}

func (r *raftLog) StoreLog(log *raft.Log) error {
	//TODO implement me
	panic("implement me")
}

func (r *raftLog) StoreLogs(logs []*raft.Log) error {
	//TODO implement me
	panic("implement me")
}

func (r *raftLog) DeleteRange(min, max uint64) error {
	//TODO implement me
	panic("implement me")
}
