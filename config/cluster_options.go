package config

import "time"

// Config holds configuration options for a distributed database.
type Config struct {
	// ReplicationFactor specifies the number of replicas for each piece of data,
	// impacting redundancy and availability.
	ReplicationFactor int

	// ShardingStrategy defines the strategy for data sharding, determining how data
	// is partitioned and distributed across nodes in the cluster.
	ShardingStrategy string

	// SchedulingStrategy specifies the task scheduling strategy, affecting data balancing
	// and load distribution.
	SchedulingStrategy string

	// LogDataStorage specifies the storage type for log data, which could be disk,
	// network storage, etc.
	LogDataStorage string

	// LogDataStoragePath is the path for storing log data.
	LogDataStoragePath string

	// SnapshotStorage specifies the storage type for snapshot data, used for backup
	// and restoration.
	SnapshotStorage string

	// SnapshotStoragePath is the path for storing snapshot data.
	SnapshotStoragePath string

	// LogDataStorageSize specifies the maximum capacity for log data storage.
	LogDataStorageSize int64

	// HeartbeatInterval defines the interval for heartbeats, used to maintain communication
	// and state synchronization among nodes in the cluster.
	HeartbeatInterval time.Duration

	// MetaNodes contains the addresses of metadata nodes, used for managing the cluster's
	// metadata information.
	MetaNodes []string

	// StoreNodes contains the addresses of storage nodes, used for storing and
	// accessing actual data.
	StoreNodes []string
}

// RegionConfig encapsulates configuration and boundary information for a specific region
// within a distributed Raft-based database.
type RegionConfig struct {
	// Options contains a set of configuration options specific to the behavior of the Raft region.
	Options Options

	// Config holds additional configuration settings related to the operation of the Raft region.
	Config Config

	// Id represents the unique identifier for the Raft region.
	Id int64

	// Start specifies the starting boundary key of the Raft region.
	Start []byte

	// End specifies the ending boundary key of the Raft region.
	End []byte
}

// StoreConfig encapsulates configuration and identification information for a store
// within a distributed system.
type StoreConfig struct {
	// Options contains a set of configuration options specific to the behavior of the store.
	Options Options

	// Config holds additional configuration settings related to the operation of the store.
	Config Config

	// Id represents the unique identifier for the store.
	Id int64

	// Addr specifies the network address at which the store can be accessed.
	Addr string
}
