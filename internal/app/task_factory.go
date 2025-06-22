package app

import (
	"fmt"
	"github.com/devenairevo/task-management/internal/contracts/tasker"
	"github.com/devenairevo/task-management/internal/task"
)

func NewTaskManager(taskDriver string) (tasker.Tasker, error) {
	switch taskDriver {
	case "local":
		return task.NewLocalTaskManager(), nil
	default:
		return nil, fmt.Errorf("unsupported task driver: %s", taskDriver)
	}
}
