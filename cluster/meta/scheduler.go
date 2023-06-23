package meta

import (
	"fmt"
	"time"
)

type NodeStatus struct {
	ID            string    // 节点唯一标识符
	Address       string    // 节点网络地址
	Alive         bool      // 节点是否存活
	Load          float64   // 节点负载信息
	LastHeartbeat time.Time // 最后一次心跳时间
	Capacity      int       // 节点容量信息
	Metadata      string    // 节点元数据
	Availability  bool      // 节点可用性状态
}

type Scheduler struct {
	config   *SchedulerConfig   // 调度器配置
	strategy SchedulingStrategy // 调度策略
}

type SchedulerConfig struct {
	MaxLoad              int    // 最大负载阈值
	MinLoad              int    // 最小负载阈值
	ReplicationFactor    int    // 副本因子
	DataShardingStrategy string // 数据分片策略
	PriorityStrategy     string // 优先级策略
	SchedulingAlgorithm  string // 任务调度算法
}

type SchedulingStrategy interface {
	SelectNode(nodes []*NodeStatus, key []byte) (*NodeStatus, error)
}

type RoundRobinScheduler struct {
	nodes []*NodeStatus
	next  int // 下一个节点索引
}

func NewScheduler(config *SchedulerConfig) *Scheduler {
	return &Scheduler{
		config:   config,
		strategy: nil, // 调度策略在设置时进行初始化
	}
}

func (s *RoundRobinScheduler) SelectNode(key []byte) (*NodeStatus, error) {
	if len(s.nodes) == 0 {
		return nil, fmt.Errorf("no available nodes")
	}

	node := s.nodes[s.next]
	s.next = (s.next + 1) % len(s.nodes)
	return node, nil
}

func (s *Scheduler) SetSchedulingStrategy(strategy SchedulingStrategy) {
	s.strategy = strategy
}

func (s *Scheduler) ScheduleNode(key []byte) (*NodeStatus, error) {
	panic("")
}
