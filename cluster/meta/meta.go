package meta

import (
	"github.com/ByteStorage/FlyDB/cluster/region"
	"github.com/ByteStorage/FlyDB/cluster/store"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/dirtree"
	"sync"
	"time"
)

type Meta struct {
	clusterConfig *config.Config
	heartbeat     map[string]time.Time
	dirTree       *dirtree.DirTree
	scheduler     *Scheduler
}

type MetadataManager interface {
	GetStore(addr string) (*store.Store, error)
	AddStore(addr string) error
	RemoveStore(addr string) error
	GetAllStores() []*store.Store
	GetRegionByID(id uint64) (*region.Region, error)
	GetStoreByID(id uint64) (*store.Store, error)
}

type metadataManager struct {
	heartbeat map[string]time.Time
	dirTree   *dirtree.DirTree
	stores    map[string]*store.Store
	regions   map[uint64]*region.Region
	mu        sync.RWMutex
}
