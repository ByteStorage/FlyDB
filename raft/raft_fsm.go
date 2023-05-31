package master

import (
	"github.com/hashicorp/raft"
	"io"
)

var _ raft.FSM = &Cluster{}

// Apply applies a Raft log entry to the FSM.
func (c *Cluster) Apply(log *raft.Log) interface{} {
	//TODO implement me
	panic("implement me")
}

// Snapshot returns a snapshot of the FSM.
func (c *Cluster) Snapshot() (raft.FSMSnapshot, error) {
	//TODO implement me
	panic("implement me")
}

// Restore stores the FSM from a snapshot.
func (c *Cluster) Restore(snapshot io.ReadCloser) error {
	//TODO implement me
	panic("implement me")
}
