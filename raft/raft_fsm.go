package raft

import (
	"github.com/hashicorp/raft"
	"io"
)

var _ raft.FSM = &Service{}

// Apply applies a Raft log entry to the FSM.
func (s Service) Apply(log *raft.Log) interface{} {
	//TODO implement me
	panic("implement me")
}

// Snapshot returns a snapshot of the FSM.
func (s Service) Snapshot() (raft.FSMSnapshot, error) {
	//TODO implement me
	panic("implement me")
}

// Restore stores the FSM from a snapshot.
func (s Service) Restore(snapshot io.ReadCloser) error {
	//TODO implement me
	panic("implement me")
}
