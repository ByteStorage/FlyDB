package cluster

type Server struct {
}

type Config struct {
	ReplicationFactor int
	ShardingStrategy  string
}

type ConfigChange struct {
	// 配置变更信息...
}
