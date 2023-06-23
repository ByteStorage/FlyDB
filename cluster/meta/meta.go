package meta

import (
	cluster2 "github.com/ByteStorage/FlyDB/cluster"
	"github.com/ByteStorage/FlyDB/cluster/region"
	"github.com/ByteStorage/FlyDB/cluster/store"
	"github.com/ByteStorage/FlyDB/lib/dirtree"
	"github.com/hashicorp/raft"
	"time"
)

type Meta struct {
	clusterConfig *cluster2.Config
	heartbeat     map[string]time.Time
	dirTree       *dirtree.DirTree
	scheduler     *Scheduler
	notification  chan cluster2.ConfigChange
}

type MetadataManager interface {
	GetStore(addr string) (*store.Store, error)
	AddStore(addr string) error
	RemoveStore(addr string) error
	GetAllStores() []*store.Store
}

type metadataManager struct {
	heartbeat  map[string]time.Time
	dirTree    *dirtree.DirTree
	stores     map[string]*store.Store
	regions    map[uint64]*region.Region
	raftGroups map[uint64]*raft.Raft
}
