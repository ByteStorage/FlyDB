package meta

import (
	"sync"
	"time"
)

type NodeStatus struct {
	ID            string    // node id
	Address       string    // node address
	Alive         bool      // whether the node is alive
	LastHeartbeat time.Time // last heartbeat time
	Capacity      int       // node capacity
	Availability  bool      // whether the node is available
}

type Scheduler struct {
	config   *SchedulerConfig   // scheduler config
	strategy SchedulingStrategy // scheduling strategy
	mu       sync.RWMutex       // mutex, to protect the scheduler
}

type SchedulerConfig struct {
	MaxLoad              int    // max load of a node
	MinLoad              int    // min load of a node
	ReplicationFactor    int    // replication factor
	DataShardingStrategy string // data sharding strategy
	PriorityStrategy     string // priority strategy
	SchedulingAlgorithm  string // scheduling algorithm
}

// SchedulingStrategy is the interface that wraps the basic SelectNode method.
type SchedulingStrategy interface {
	// SelectNode selects num node to schedule the key.
	SelectNode(nodes []*NodeStatus, key []byte, num int) (*NodeStatus, error)
}

// SetSchedulingStrategy sets the scheduling strategy.
func (s *Scheduler) SetSchedulingStrategy(strategy SchedulingStrategy) {
	s.strategy = strategy
}

// ScheduleNode schedules the key to a node.
func (s *Scheduler) ScheduleNode(key []byte) (*NodeStatus, error) {
	panic("implement me")
}
