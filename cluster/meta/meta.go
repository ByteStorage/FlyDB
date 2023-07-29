package meta

import (
	"fmt"
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
	// ApplyConfig applies a new config to the cluster.
	ApplyConfig(config *config.Config) error
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

// GetStore gets a store by address.
func (m *meta) GetStore(addr string) (*store.Store, error) {
	//TODO implement me
	panic("implement me")
}

// AddStore adds a new store to the cluster.
func (m *meta) AddStore(addr string) error {
	//TODO implement me
	panic("implement me")
}

// RemoveStore removes a store from the cluster.
func (m *meta) RemoveStore(addr string) error {
	//TODO implement me
	panic("implement me")
}

// GetAllStores gets all stores in the cluster.
func (m *meta) GetAllStores() []*store.Store {
	//TODO implement me
	panic("implement me")
}

// GetRegionByID gets a region by id.
func (m *meta) GetRegionByID(id uint64) (*region.Region, error) {
	//TODO implement me
	panic("implement me")
}

// GetStoreByID gets a store by id.
func (m *meta) GetStoreByID(id uint64) (*store.Store, error) {
	//TODO implement me
	panic("implement me")
}

// ApplyConfig applies a new config to the cluster.
func (m *meta) ApplyConfig(config *config.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m != nil {
		err := m.stop()
		if err != nil {
			// if err = ErrNotStarted, it means the meta node has not started yet.
		}
		err = m.start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *meta) start() error {
	for _, metaNode := range m.clusterConfig.MetaNodes {
		// ssh to the meta node
		// start the meta node
		fmt.Println(metaNode)
		panic("implement me")
	}
	return nil
}

func (m *meta) stop() error {
	//TODO implement me
	panic("implement me")
}

// NewMeta creates a new meta.
func NewMeta(conf config.Config) MetadataManager {
	return &meta{
		clusterConfig: &conf,
		heartbeat:     make(map[string]time.Time),
		dirTree:       dirtree.NewDirTree(),
		stores:        make(map[string]*store.Store),
		regions:       make(map[uint64]*region.Region),
	}
}
