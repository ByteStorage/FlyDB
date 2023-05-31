package master

import "github.com/hashicorp/raft"

// Interface is the interface that must be implemented by a raft server.
type Interface interface {
	State() raft.RaftState
	Peers() ([]string, error)
	Close() error
	IsLeader() bool
	IsCandidate() bool
	Leader() string
	Apply(b []byte) error
	AddServer(addr string) error
	ShowDebugInfo(witch string) ([]byte, error)
	UserSnapshot() error
}
