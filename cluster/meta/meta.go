package meta

import (
	"github.com/ByteStorage/FlyDB/cluster/store"
	"time"
)

type Meta struct {
	store     map[string]*store.Store
	heartbeat map[string]time.Time
}

type Manager interface {
	// GetStore gets store by addr
	GetStore(addr string) (*store.Store, error)
}
