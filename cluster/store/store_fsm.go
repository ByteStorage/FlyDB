package store

import (
	"github.com/hashicorp/raft"
	"io"
)

// fsm implements raft.FSM interface
type fsm struct {
	//implement me
}

func newFSM() raft.FSM {
	return &fsm{}
}

func (f fsm) Apply(log *raft.Log) interface{} {
	//TODO implement me
	panic("implement me")
}

func (f fsm) Snapshot() (raft.FSMSnapshot, error) {
	//TODO implement me
	panic("implement me")
}

func (f fsm) Restore(snapshot io.ReadCloser) error {
	//TODO implement me
	panic("implement me")
}
