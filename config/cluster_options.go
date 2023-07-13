package config

import "time"

type Config struct {
	ReplicationFactor  int
	ShardingStrategy   string
	SchedulingStrategy string
	LogDataStorage     string
	LogDataStoragePath string
	LogDataStorageSize int64
	HeartbeatInterval  time.Duration
}
