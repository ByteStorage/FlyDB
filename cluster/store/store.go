package store

import (
	"github.com/ByteStorage/FlyDB/cluster/region"
	"sync"
)

// store stores the data of a store.
type store struct {
	id         uint64                    // store id
	addr       string                    // store address
	regionList map[uint64]*region.Region // region list, to store the regions in the store.
	size       int64                     // size
	mu         sync.RWMutex              // mutex, to protect the store.
}

// Store is the interface of store.
type Store interface {
	// GetRegionByKey gets region and leader peer by region key from cluster.
	GetRegionByKey(key []byte) (*region.Region, error)
	// GetRegionByID gets region and leader peer by region id from cluster.
	GetRegionByID(id uint64) (*region.Region, error)
	// AddRegion adds a new region to cluster.
	AddRegion(region *region.Region) error
	// RemoveRegion removes a region from cluster.
	RemoveRegion(id uint64) error
	// Split splits the region into two regions.
	Split(region *region.Region, splitKey []byte) error
	// Merge merges two adjacent regions into one region.
	Merge(regionA *region.Region, regionB *region.Region) error
	// GetSize gets the total size of the store.
	GetSize() int64
}
