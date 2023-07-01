package meta

import (
	"sync"
	"time"
)

type NodeStatus struct {
	ID            string
	Address       string
	Alive         bool
	LastHeartbeat time.Time
	Capacity      int
	Availability  bool
}

type Scheduler struct {
	config   *SchedulerConfig
	strategy SchedulingStrategy
	mu       sync.RWMutex
}

type SchedulerConfig struct {
	MaxLoad              int
	MinLoad              int
	ReplicationFactor    int
	DataShardingStrategy string
	PriorityStrategy     string
	SchedulingAlgorithm  string
}

type SchedulingStrategy interface {
	SelectNode(nodes []*NodeStatus, key []byte) (*NodeStatus, error)
}

func (s *Scheduler) SetSchedulingStrategy(strategy SchedulingStrategy) {
	s.strategy = strategy
}

func (s *Scheduler) ScheduleNode(key []byte) (*NodeStatus, error) {
	panic("")
}
