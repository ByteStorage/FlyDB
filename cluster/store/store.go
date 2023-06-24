package store

import "github.com/ByteStorage/FlyDB/cluster/region"

type store struct {
	addr       string
	regionList map[uint64]*region.Region
	size       int
}

type Store interface {
	GetRegion(key []byte) (*region.Region, error)
}
