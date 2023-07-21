package store

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/hashicorp/raft"
)

func newStableLog(conf config.Config) (raft.StableStore, error) {
	_ = Init()

	// Get the "memory" DataStoreFactory from the map
	return getDataStore(conf)
}
