package config

import "time"

type Config struct {
	ReplicationFactor  int
	ShardingStrategy   string
	SchedulingStrategy string
	HeartbeatInterval  time.Duration
}
