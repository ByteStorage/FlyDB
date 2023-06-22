package store

import "github.com/ByteStorage/FlyDB/cluster/region"

type Store struct {
	Addr       string
	regionList map[uint64]*region.Region
	size       int
}

type Manager interface {
	GetRegion(key []byte) (*region.Region, error)
}
