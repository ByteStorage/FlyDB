package region

import (
	"errors"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine/engine"
	"github.com/hashicorp/raft"
	"sync"
)

// The region is responsible for maintaining the key range and the raftGroups it belongs to,
// directly interacting with the underlying DB APIs,
// and the actual groups where the data falls into the disk.
// raftGroups contains the other raft node of the region, a region has at least three replicas.
// region and replicas are a raft group, and the one region is the leader of the raft group.
// region stores the data of a region.
type region struct {
	id         int64                 // region id
	startKey   []byte                // start key
	endKey     []byte                // end key
	db         *engine.DB            // db, to store the data.
	raft       *raft.Raft            // raft, to store the raft group.
	raftGroups map[uint64]*raft.Raft // raft groups, to store the raft groups of the region.
	leader     string                // leader peer
	peers      []string              // peers
	size       int64                 // size
	mu         sync.RWMutex          // mutex, to protect the region.
	conf       config.Config         // config
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
	// GetID gets the id of the region.
	GetID() int64
}

func NewRegion(conf config.RegionConfig) (Region, error) {
	db, err := engine.NewDB(conf.Options)
	if err != nil {
		return nil, errors.New("new db failed")
	}
	raftNode, err := newRaftNode(conf.Config)
	if err != nil {
		return nil, errors.New("new raft node failed")
	}
	return &region{
		startKey:   conf.Start,
		endKey:     conf.End,
		raftGroups: make(map[uint64]*raft.Raft),
		db:         db,
		mu:         sync.RWMutex{},
		conf:       conf.Config,
		raft:       raftNode,
	}, nil
}

// newRaftNode creates a new raft node for the store.
func newRaftNode(conf config.Config) (*raft.Raft, error) {
	// All new methods below can add other return values as needed, such as err

	// create default config for raft
	raftConfig := newDefaultConfig()

	// setup Raft communication
	t := newTransport()

	// create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := newSnapshotStore(conf)
	if err != nil {
		return nil, err
	}

	// create the log store and stable store
	logStore, err := newRaftLog(conf)
	if err != nil {
		return nil, err
	}
	stableStore, err := newStableLog(conf)
	if err != nil {
		return nil, err
	}

	// create a new finite state machine
	f := newFSM()

	// instantiate the Raft system
	r, err := raft.NewRaft(raftConfig, f, logStore, stableStore, snapshots, t)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *region) Put(key []byte, value []byte) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.db.Put(key, value)
}

func (r *region) Get(key []byte) ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.db.Get(key)
}

func (r *region) Delete(key []byte) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.db.Delete(key)
}

func (r *region) GetStartKey() []byte {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.startKey
}

func (r *region) GetEndKey() []byte {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.endKey
}

func (r *region) GetLeader() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.leader
}

func (r *region) GetPeers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.peers
}

func (r *region) TransferLeader(peer string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.peers {
		if p == peer {
			err := r.raft.LeadershipTransferToServer(raft.ServerID(peer), raft.ServerAddress(peer)).Error()
			if err != nil {
				return err
			}
			r.leader = peer
			return nil
		}
	}
	return errors.New("no such peer exists")
}

func (r *region) AddPeer(peer string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	contains := func(arr []string, it string) bool {
		for _, n := range arr {
			if n == it {
				return true
			}
		}
		return false
	}
	if !contains(r.peers, peer) {
		r.peers = append(r.peers, peer)
		return r.raft.AddVoter(raft.ServerID(peer), raft.ServerAddress(peer), 0, 0).Error()
	}
	return errors.New("peer already exists")
}

func (r *region) RemovePeer(peer string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := 0; i < len(r.peers); i++ {
		if r.peers[i] == peer {
			r.peers = append(r.peers[:i], r.peers[i+1:]...)
			return r.raft.DemoteVoter(raft.ServerID(peer), 0, 0).Error()
		}
	}
	return errors.New("the specified peer does not exist")
}

func (r *region) GetSize() int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.size
}

func (r *region) GetID() int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.id
}
