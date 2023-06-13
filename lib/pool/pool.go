package pool

type Task struct {
	// type of task
	Type int    // 0: read, 1: write, 2: delete, 3: update
	Key  []byte // key of task
	Val  []byte // value of task
}

// Pool is a connection pool interface
type Pool interface {
	// GetTask get a task from pool
	GetTask() (*Task, error)
	// PutTask put a task into pool
	PutTask(task *Task) error
	// GetTaskNum get the number of tasks in the pool
	GetTaskNum() int
}
