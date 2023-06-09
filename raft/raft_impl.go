package cluster

import (
	"github.com/hashicorp/raft"
)

var _ Interface = &Impl{}

// Impl is the interface that must be implemented by a raft server.
type Impl struct {
	m Master
}

func (i *Impl) State() raft.RaftState {
	return i.m.Raft.State()
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
	return i.m.Raft.State() == raft.Leader
}

func (i *Impl) IsCandidate() bool {
	return i.m.Raft.State() == raft.Candidate
}

func (i *Impl) Leader() string {
	address, _ := i.m.Raft.LeaderWithID()
	return string(address)
}

func (i *Impl) Apply(b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) AddServer(addr string) error {
	voter := i.m.Raft.AddVoter(raft.ServerID(addr), raft.ServerAddress(addr), 0, 0)
	if voter.Error() != nil {
		return voter.Error()
	}
	return nil
}

func (i *Impl) ShowDebugInfo(witch string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) UserSnapshot() error {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) Put(key, value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) Get(key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
