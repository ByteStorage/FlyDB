package store

import (
	"github.com/hashicorp/raft"
	"io"
)

// snapshot implements raft.SnapshotStore interface
type snapshot struct {
	//implement me
}

func newSnapshot() raft.SnapshotStore {
	return &snapshot{}
}

func (s snapshot) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
	//TODO implement me
	panic("implement me")
}

func (s snapshot) List() ([]*raft.SnapshotMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (s snapshot) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}
