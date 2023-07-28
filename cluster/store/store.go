package store

import (
	"bytes"
	"errors"
	"github.com/ByteStorage/FlyDB/cluster/region"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/bwmarrin/snowflake"
	"sync"
)

var (
	MinKey []byte                                           // Min Key of the all regions
	MaxKey = []byte{255, 255, 255, 255, 255, 255, 255, 255} // Max Key of the all regions
)

// The store component is responsible for managing the division and merging of region partitions.
// All regions under the store share a port number.
// Each region under the store is in a raftGroups, and the region clusters in the raftGroups communicate through grpc
// store stores the data of a store.
type store struct {
	id         int64 // store id
	conf       config.Config
	opts       config.Options
	addr       string                  // store address
	regionList map[int64]region.Region // region list, to store the regions in the store.
	mu         sync.RWMutex            // mutex, to protect the store.
	node       *snowflake.Node         // snowflake node, to generate the id.
}

// Store is the interface of store.
type Store interface {
	// GetRegionByKey gets region and leader peer by region key from cluster.
	GetRegionByKey(key []byte) (region.Region, error)
	// GetRegionByID gets region and leader peer by region id from cluster.
	GetRegionByID(id int64) (region.Region, error)
	// Split splits the region into two regions.
	Split(region region.Region, splitKey []byte) error
	// Merge merges two adjacent regions into one region.
	Merge(regionA region.Region, regionB region.Region) error
	// GetSize gets the total size of the store.
	GetSize() int64
}

// NewStore creates a new store.
func NewStore(conf config.StoreConfig) (Store, error) {
	// create a new region, when initialize, a store just has one region.
	// when the region size exceeds the threshold, the region will be split into two regions.
	newRegion, err := region.NewRegion(MinKey, MaxKey, conf.Options, conf.Config)
	if err != nil {
		return nil, err
	}
	// create a new snowflake node.
	node, err := snowflake.NewNode(conf.Id)
	if err != nil {
		return nil, err
	}
	return &store{
		id:   conf.Id,
		node: node,
		regionList: map[int64]region.Region{
			newRegion.GetID(): newRegion,
		},
		addr: conf.Addr,
		conf: conf.Config,
		opts: conf.Options,
		mu:   sync.RWMutex{},
	}, nil
}

func (s *store) GetRegionByKey(key []byte) (region.Region, error) {
	for _, r := range s.regionList {
		if isKeyInRange(key, r.GetStartKey(), r.GetEndKey()) {
			return r, nil
		}
	}
	return nil, errors.New("the specified region does not exist")
}

func (s *store) GetRegionByID(id int64) (region.Region, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.regionList[id]; !ok {
		return nil, errors.New("the specified region does not exist")
	}
	return s.regionList[id], nil
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
