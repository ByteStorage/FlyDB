package meta

import (
	"fmt"
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

type RoundRobinScheduler struct {
	nodes []*NodeStatus
	next  int
}

func NewScheduler(config *SchedulerConfig) *Scheduler {
	return &Scheduler{
		config:   config,
		strategy: nil,
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
