package queue

import (
	"github.com/devenairevo/task-management/internal/contracts/queuer"
	"github.com/devenairevo/task-management/internal/contracts/tasker"
	"github.com/devenairevo/task-management/internal/types"
)

func StartQueueWorker(channel queuer.Queuer, taskMng tasker.Tasker) {
	go func(channel queuer.Queuer, taskMng tasker.Tasker) {
		for {
			finishedTask, err := channel.Dequeue()
			if err != nil {
				continue
			}

			t, err := taskMng.DescribeTask(finishedTask.ID)
			if err == nil {
				t.Status = types.Done
			}
		}
	}(channel, taskMng)
}
