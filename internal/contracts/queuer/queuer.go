package queuer

import (
	"github.com/devenairevo/task-management/internal/task"
)

type Queuer interface {
	Enqueue(*task.Task) error
	Dequeue() (*task.Task, error)
	IsEmpty() bool
	Size() int
}
