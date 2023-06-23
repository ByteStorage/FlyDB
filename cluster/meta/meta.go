package meta

import (
	"github.com/ByteStorage/FlyDB/cluster/store"
	"github.com/ByteStorage/FlyDB/lib/dirtree"
	"time"
)

type Meta struct {
	store     map[string]*store.Store
	heartbeat map[string]time.Time
	dirTree   *dirtree.DirTree
}

type Manager interface {
	// GetStore gets store by addr
	GetStore(addr string) (*store.Store, error)
}
