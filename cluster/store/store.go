package store

import (
	"bytes"
	"errors"
	"github.com/ByteStorage/FlyDB/cluster/region"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/hashicorp/raft"
	"sync"
)

// The store component is responsible for managing the division and merging of region partitions.
// All regions under the store share a port number.
// Each region under the store is in a raftGroups, and the region clusters in the raftGroups communicate through grpc
// store stores the data of a store.
type store struct {
	id         string // store id
	conf       config.Config
	opts       config.Options
	addr       string                   // store address
	regionList map[uint64]region.Region // region list, to store the regions in the store.
	size       int64                    // size
	mu         sync.RWMutex             // mutex, to protect the store.
	raft       *raft.Raft               // raft, to store the raft group.
}

// Store is the interface of store.
type Store interface {
	// GetRegionByKey gets region and leader peer by region key from cluster.
	GetRegionByKey(key []byte) (region.Region, error)
	// GetRegionByID gets region and leader peer by region id from cluster.
	GetRegionByID(id uint64) (region.Region, error)
	// AddRegion adds a new region to cluster.
	AddRegion(region region.Region) error
	// RemoveRegion removes a region from cluster.
	RemoveRegion(id uint64) error
	// Split splits the region into two regions.
	Split(region region.Region, splitKey []byte) error
	// Merge merges two adjacent regions into one region.
	Merge(regionA region.Region, regionB region.Region) error
	// GetSize gets the total size of the store.
	GetSize() int64
}

func (s *store) GetRegionByKey(key []byte) (region.Region, error) {
	for _, r := range s.regionList {
		if isKeyInRange(key, r.GetStartKey(), r.GetEndKey()) {
			return r, nil
		}
	}
	return nil, errors.New("the specified region does not exist")
}

func (s *store) GetRegionByID(id uint64) (region.Region, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.regionList[id]; !ok {
		return nil, errors.New("the specified region does not exist")
	}
	return s.regionList[id], nil
}

func (s *store) AddRegion(region region.Region) error {
	panic("implement me")
}

func (s *store) RemoveRegion(id uint64) error {
	panic("implement me")
}

func (s *store) Split(region region.Region, splitKey []byte) error {
	panic("implement me")
}

func (s *store) Merge(regionA region.Region, regionB region.Region) error {
	panic("implement me")
}

func (s *store) GetSize() int64 {
	panic("implement me")
}

// isKeyInRange checks if the key is in the range of the region.
func isKeyInRange(key, startRange, endRange []byte) bool {
	// Compare the key to the start of the range
	// If key < startRange, it's not in range
	if bytes.Compare(key, startRange) < 0 {
		return false
	}

	// Compare the key to the end of the range
	// If key >= endRange, it's not in range
	if bytes.Compare(key, endRange) >= 0 {
		return false
	}

	// If neither of the above, the key is in range
	return true
}
