package region

import (
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/hashicorp/raft"
)

type Region struct {
	id       uint64
	startKey []byte
	endKey   []byte
	db       *engine.DB
	raft     *raft.Raft
}

type Manager interface {
	// GetRegionByKey gets region and leader peer by region key from cluster.
	GetRegionByKey(key []byte) (*Region, error)
	// GetRegionByID gets region and leader peer by region id from cluster.
	GetRegionByID(id uint64) (*Region, error)
	// Put puts a key-value pair to region.
	Put(key []byte, value []byte) error
	// Get gets value by key from region.
	Get(key []byte) ([]byte, error)
	// Delete deletes a key-value pair from region.
	Delete(key []byte) error
}
