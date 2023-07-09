package meta

import (
	"errors"
	"github.com/ByteStorage/FlyDB/cluster/region"
	"github.com/ByteStorage/FlyDB/cluster/store"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/dirtree"
	"github.com/hashicorp/raft"
	"sync"
	"time"
)

// MetadataManager manages the metadata of the cluster.
type MetadataManager interface {
	// GetStore gets a store by address.
	GetStore(addr string) (*store.Store, error)
	// AddStore adds a new store to the cluster.
	AddStore(addr string) error
	// RemoveStore removes a store from the cluster.
	RemoveStore(addr string) error
	// GetAllStores gets all stores in the cluster.
	GetAllStores() []*store.Store
	// GetRegionByID gets a region by id.
	GetRegionByID(id uint64) (*region.Region, error)
	// GetStoreByID gets a store by id.
	GetStoreByID(id uint64) (*store.Store, error)
}

// meta stores the metadata of the cluster.
// meta will manage all stores and regions in the cluster.
// meta has at least three nodes in the cluster.
// meta nodes will create a raft group to manage the metadata of the cluster.
type meta struct {
	clusterConfig *config.Config            // cluster config, including cluster id, cluster name, etc.
	heartbeat     map[string]time.Time      // stores heartbeat, to check whether a store is alive.
	dirTree       *dirtree.DirTree          // directory tree, to store data directory.
	stores        map[string]*store.Store   // stores, to store the stores in the cluster.
	regions       map[uint64]*region.Region // regions, to store the regions in the cluster.
	mu            sync.RWMutex              // mutex, to protect the metadata.
	scheduler     *Scheduler                // scheduler, to schedule the cluster.
	raft          *raft.Raft                // raft, to store the raft group.
}

func (m *meta) GetStore(addr string) (*store.Store, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	store, ok := m.stores[addr]
	if !ok {
		return nil, errors.New("store not found")
	}

	return store, nil
}

func (m *meta) AddStore(addr string) error {
	panic("implement me")
}

func (m *meta) RemoveStore(addr string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.stores[addr]; !ok {
		return errors.New("store not found")
	}

	// Remove the store from the metadata
	delete(m.stores, addr)

	return nil
}

func (m *meta) GetAllStores() []*store.Store {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stores := make([]*store.Store, 0, len(m.stores))
	for _, store := range m.stores {
		stores = append(stores, store)
	}

	return stores
}

func (m *meta) GetRegionByID(id uint64) (*region.Region, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	region, ok := m.regions[id]
	if !ok {
		return nil, errors.New("region not found")
	}

	return region, nil
}

func (m *meta) GetStoreByID(id uint64) (*store.Store, error) {
	panic("implement me")
}
