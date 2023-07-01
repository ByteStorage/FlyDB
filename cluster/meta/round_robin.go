package meta

import "fmt"

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
