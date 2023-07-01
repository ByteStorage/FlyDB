package meta

// RoundRobinScheduler is a scheduler that selects nodes in round-robin order.
type RoundRobinScheduler struct {
	nodes []*NodeStatus // all nodes
	next  int           // next node index
}

// NewRoundRobinScheduler creates a new RoundRobinScheduler.
func NewRoundRobinScheduler(config *SchedulerConfig) *Scheduler {
	return &Scheduler{
		config:   config,
		strategy: nil,
	}
}

// SelectNode selects a node in round-robin order.
func (r *RoundRobinScheduler) SelectNode(nodes []*NodeStatus, key []byte, num int) (*NodeStatus, error) {
	panic("implement me")
}
