package cluster

import (
	"github.com/ByteStorage/FlyDB/cluster/meta"
	"github.com/ByteStorage/FlyDB/config"
)

// server is the implementation of Cluster interface.
type server struct {
	meta *meta.MetadataManager // metadata manager, to manage the metadata of the cluster.
}

// Cluster is the interface of cluster.
type Cluster interface {
	// Put puts a key-value pair to cluster.
	Put(key, value []byte) error
	// Get gets value by key from cluster.
	Get(key []byte) ([]byte, error)
	// Delete deletes a key-value pair from cluster.
	Delete(key []byte) error
	// GetKeys gets all keys from cluster.
	GetKeys() ([][]byte, error)
	// UpdateConfig updates the config of cluster.
	UpdateConfig(config *config.Config) error
}

// NewCluster creates a new cluster.
func NewCluster(config *config.Config) (Cluster, error) {
	panic("implement me")
}
