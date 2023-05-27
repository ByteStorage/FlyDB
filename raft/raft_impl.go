package raft

import (
	"github.com/hashicorp/raft"
)

var _ Interface = &Impl{}

// Impl is the interface that must be implemented by a raft server.
type Impl struct {
	c Cluster
}

func (i *Impl) State() raft.RaftState {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) Peers() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) Close() error {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) IsLeader() bool {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) IsCandidate() bool {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) Leader() string {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) Apply(b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) AddServer(addr string) error {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) ShowDebugInfo(witch string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) UserSnapshot() error {
	//TODO implement me
	panic("implement me")
}
