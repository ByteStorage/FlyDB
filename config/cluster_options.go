package config

import (
	"github.com/hashicorp/raft"
	"time"
)

type Config struct {
	ReplicationFactor   int
	ShardingStrategy    string
	SchedulingStrategy  string
	LogDataStorage      string
	LogDataStoragePath  string
	SnapshotStorage     string
	SnapshotStoragePath string
	LogDataStorageSize  int64
	HeartbeatInterval   time.Duration
	HeartbeatTimeout    time.Duration
	LocalAddress        raft.ServerAddress
}
