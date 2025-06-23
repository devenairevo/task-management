package app

import (
	"fmt"
	"github.com/devenairevo/task-management/internal/contracts/queuer"
	"github.com/devenairevo/task-management/internal/queue"
	"sync"
)

func NewQueueManager(buffSize int, queueDriver string, wg *sync.WaitGroup) (queuer.Queuer, error) {
	switch queueDriver {
	case "local":
		return queue.NewChannel(buffSize, wg), nil
	default:
		return nil, fmt.Errorf("unsupported queue driver: %s", queueDriver)
	}
}
