package master

import (
	"github.com/hashicorp/raft"
)

var _ Interface = &Impl{}

// Impl is the interface that must be implemented by a raft server.
type Impl struct {
	c Cluster
}

func (i *Impl) State() raft.RaftState {
	return i.c.Raft.State()
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
	address, _ := i.c.Raft.LeaderWithID()
	i.c.Leader = string(address)
	return string(address)
}

func (i *Impl) Apply(b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) AddServer(addr string) error {
	voter := i.c.Raft.AddVoter(raft.ServerID(addr), raft.ServerAddress(addr), 0, 0)
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
