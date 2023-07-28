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
	MinKey    []byte                                                  // Min Key of the all regions
	MaxKey           = []byte{255, 255, 255, 255, 255, 255, 255, 255} // Max Key of the all regions
	Threshold int64  = 256 * 1024 * 1024                              // Threshold of the region size
)

// The store component is responsible for managing the division and merging of region partitions.
// All regions under the store share a port number.
// Each region under the store is in a raftGroups, and the region clusters in the raftGroups communicate through grpc
// store stores the data of a store.
type store struct {
	id         int64 // store id
	conf       config.Config
	opts       config.Options
	addr       string          // store address
	regionList []region.Region // region list, to store the regions in the store.
	mu         sync.RWMutex    // mutex, to protect the store.
	node       *snowflake.Node // snowflake node, to generate the id.
}

// Store is the interface of store.
type Store interface {
	// GetRegionByKey gets region and leader peer by region key from cluster.
	GetRegionByKey(key []byte) (region.Region, error)
	// Split splits the region into two regions.
	Split() error
	// Merge merges two adjacent regions into one region.
	Merge() error
}

// NewStore creates a new store.
func NewStore(conf config.StoreConfig) (Store, error) {
	// create a new region config.
	regionConfig := config.RegionConfig{
		Options: conf.Options,
		Config:  conf.Config,
		Start:   MinKey,
		End:     MaxKey,
	}
	// create a new region, when initialize, a store just has one region.
	// when the region size exceeds the threshold, the region will be split into two regions.
	newRegion, err := region.NewRegion(regionConfig)
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
		regionList: []region.Region{
			newRegion,
		},
		addr: conf.Addr,
		conf: conf.Config,
		opts: conf.Options,
		mu:   sync.RWMutex{},
	}, nil
}

// GetRegionByKey gets region by region key from store.
func (s *store) GetRegionByKey(key []byte) (region.Region, error) {
	for _, r := range s.regionList {
		if isKeyInRange(key, r.GetStartKey(), r.GetEndKey()) {
			return r, nil
		}
	}
	return nil, errors.New("the specified region does not exist")
}

func (s *store) Split() error {
	for _, r := range s.regionList {
		if r.GetSize() >= Threshold {
			// define the middle key
			end := r.GetEndKey()
			start := r.GetStartKey()
			middle := make([]byte, len(end))
			for i := 0; i < len(end); i++ {
				middle[i] = (end[i] + start[i]) / 2
			}
			// create a new region config
			regionConfig := config.RegionConfig{
				Options: s.opts,
				Config:  s.conf,
				Start:   middle,
				End:     end,
			}
			// create a new region
			newRegion, err := region.NewRegion(regionConfig)
			if err != nil {
				return err
			}
			// move the data to the new region
			err = moveDataToNewRegion(newRegion, middle, end, r, start, end)
			if err != nil {
				return err
			}
			// add the new region to the region list
			s.regionList = append(s.regionList, newRegion)
		}
	}
	return nil
}

// Merge merges two adjacent regions into one region.
func (s *store) Merge() error {
	// find the two adjacent regions which have the smallest size
	for i, r := range s.regionList {
		if i == len(s.regionList)-1 {
			break
		}
		// regionList[i] and regionList[i+1] are adjacent
		// and if the size of them is smaller than the threshold / 2, merge them
		if r.GetSize()+s.regionList[i+1].GetSize() < Threshold/2 {
			err := mergeTwoRegions(r, s.regionList[i+1])
			if err != nil {
				return err
			}
			// delete the second region
			s.regionList = append(s.regionList[:i+1], s.regionList[i+2:]...)
		}
	}
	return nil
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

// moveDataToNewRegion moves the data from the old region to the new region.
// new: the new region
// newStartKey: the start key of the new region
// newEndKey: the end key of the new region
// old: the old region
// oldStartKey: the start key of the old region
// oldEndKey: the end key of the old region
func moveDataToNewRegion(new region.Region, newStartKey []byte, newEndKey []byte, old region.Region, oldStartKey []byte, oldEndKey []byte) error {
	// modify the start key and the end key of the old region
	panic("implement me")
}

// mergeTwoRegions merges two adjacent regions into one region.
// r1: the first region
// r2: the second region
// move the data from the second region to the first region
func mergeTwoRegions(r1 region.Region, r2 region.Region) error {
	panic("implement me")
}
