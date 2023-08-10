package meta

import (
	"github.com/hashicorp/raft"
	"io"
)

// fsm is the finite state machine of the raft group for meta.
type fsm struct {
}

func newFSM() *fsm {
	panic("implement me")
}

func (f *fsm) Apply(log *raft.Log) interface{} {
	panic("implement me")
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	panic("implement me")
}

func (f *fsm) Restore(rc io.ReadCloser) error {
	panic("implement me")
}
