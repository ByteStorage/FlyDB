package region

import (
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/hashicorp/raft"
	"sync"
)

type region struct {
	id         uint64
	startKey   []byte
	endKey     []byte
	db         *engine.DB
	raft       *raft.Raft
	raftGroups map[uint64]*raft.Raft
	leader     string
	peers      []string
	size       int64
	mu         sync.RWMutex
}

type Region interface {
	// Put puts a key-value pair to region.
	Put(key []byte, value []byte) error
	// Get gets value by key from region.
	Get(key []byte) ([]byte, error)
	// Delete deletes a key-value pair from region.
	Delete(key []byte) error
	// GetStartKey gets the start key of the region.
	GetStartKey() []byte
	// GetEndKey gets the end key of the region.
	GetEndKey() []byte
	// GetLeader gets the leader peer of the region.
	GetLeader() string
	// GetPeers gets the peers of the region.
	GetPeers() []string
	// TransferLeader transfers the leader role to a specified peer.
	TransferLeader(peer string) error
	// AddPeer adds a new peer to the raft group.
	AddPeer(peer string) error
	// RemovePeer removes a peer from the raft group.
	RemovePeer(peer string) error
	// GetSize gets the total size of the region.
	GetSize() int64
}
