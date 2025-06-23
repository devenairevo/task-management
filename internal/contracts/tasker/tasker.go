package tasker

import (
	"github.com/devenairevo/task-management/internal/task"
)

type Tasker interface {
	CreateTask(*task.CreateTaskParams) (*task.Task, error)
	ListTasks() ([]*task.Task, error)
	DescribeTask(taskID int) (*task.Task, error)
	UpdateTask(*task.Task) error
}
