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
	return i.m.Peers, nil
}

func (i *Impl) Close() error {
	// close the raft node
	err := i.m.Raft.Shutdown().Error()
	if err != nil {
		return nil
	}
	// close the raft log store
	err = i.m.RaftLog.Close()
	return err
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

	// call the *Cluster.Apply method
	af := i.m.Raft.Apply(b, 0)

	// Error will only return generic errors related to raft
	if err := af.Error(); err != nil {
		return err
	}

	// *Cluster.Apply will return an error by Response, not by Error
	// so here we must check whether the response is an error
	response := af.Response()
	if err, ok := response.(error); ok {
		return err
	}

	return nil
}

func (i *Impl) AddServer(addr string) error {
	return i.m.Raft.AddVoter(raft.ServerID(addr), raft.ServerAddress(addr), 0, 0).Error()
}

func (i *Impl) ShowDebugInfo(witch string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) UserSnapshot() error {
	sf := i.m.Raft.Snapshot()
	return sf.Error()
}

func (i *Impl) Put(key, value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) Get(key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
