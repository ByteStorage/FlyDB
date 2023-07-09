package meta

import (
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
