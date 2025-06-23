package app

import (
	"fmt"
	"github.com/devenairevo/task-management/internal/contracts/tasker"
	"github.com/devenairevo/task-management/internal/file"
	"github.com/devenairevo/task-management/internal/task"
)

func NewFileTaskManager(taskDriver string, dirName string) (tasker.Tasker, error) {
	switch taskDriver {
	case "file":
		return file.NewManager(dirName), nil
	case "local":
		return task.NewLocalTaskManager(), nil
	default:
		return nil, fmt.Errorf("unsupported task driver: %s", taskDriver)
	}
}
