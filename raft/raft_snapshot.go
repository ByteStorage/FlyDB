package master

import "github.com/hashicorp/raft"

var _ raft.FSMSnapshot = &FSMSnapshot{}

func (F *FSMSnapshot) Persist(sink raft.SnapshotSink) error {
	//TODO implement me
	panic("implement me")
}

func (F *FSMSnapshot) Release() {
	//TODO implement me
	panic("implement me")
}
