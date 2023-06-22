package region

type Region struct {
	Id       uint64
	StartKey []byte
	EndKey   []byte
}

type Manager interface {
	// GetRegionByKey gets region and leader peer by region key from cluster.
	GetRegionByKey(key []byte) (*Region, error)
	// GetRegionByID gets region and leader peer by region id from cluster.
	GetRegionByID(id uint64) (*Region, error)
}
