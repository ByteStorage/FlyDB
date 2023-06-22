package meta

import "github.com/ByteStorage/FlyDB/cluster/store"

type Meta struct {
	store map[string]*store.Store
}

type Manager interface {
	// GetStore gets store by addr
	GetStore(addr string) (*store.Store, error)
}
