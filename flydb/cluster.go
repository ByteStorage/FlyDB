package flydb

import (
	"github.com/ByteStorage/flydb/config"
	"github.com/ByteStorage/flydb/engine"
	cluster "github.com/ByteStorage/flydb/raft"
)

// NewFlyDbCluster create a new db cluster
func NewFlyDbCluster(masterList []string, slaveList []string, options config.Options) (*engine.DB, error) {
	cluster.NewRaftCluster(masterList, slaveList)
	panic("implement me")
}
