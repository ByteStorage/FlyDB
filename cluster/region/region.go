package region

import (
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/hashicorp/raft"
	"sync"
)

// region stores the data of a region.
type region struct {
	id         uint64                // region id
	startKey   []byte                // start key
	endKey     []byte                // end key
	db         *engine.DB            // db, to store the data.
	raft       *raft.Raft            // raft, to store the raft group.
	raftGroups map[uint64]*raft.Raft // raft groups, to store the raft groups of the region.
	leader     string                // leader peer
	peers      []string              // peers
	size       int64                 // size
	mu         sync.RWMutex          // mutex, to protect the region.
}

// Region is the interface of region.
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
