package queue

import (
	"errors"
	"fmt"
	"github.com/devenairevo/task-management/internal/task"
	"github.com/devenairevo/task-management/internal/types"
	"sync"
)

type Channel struct {
	mu       sync.Mutex
	wg       *sync.WaitGroup
	BuffSize int
	Task     *task.Task
	chanel   chan *task.Task
}

func NewChannel(buffSize int, wg *sync.WaitGroup) *Channel {
	return &Channel{
		BuffSize: buffSize,
		chanel:   make(chan *task.Task, buffSize),
		wg:       wg,
	}
}

func (c *Channel) Enqueue(task *task.Task) error {
	if task.ID <= 0 || task.Name == "" {
		return errors.New("error with adding task to the queue")
	}
	if c.chanel == nil {
		return errors.New("channel not created yet")
	}

	c.wg.Add(1)
	c.chanel <- task

	task.Status = types.Created
	fmt.Printf("Task with ID %d and name %s created\n", task.ID, task.Name)

	task.Status = types.Processing
	fmt.Printf("Task with ID %d and name %s started processing\n", task.ID, task.Name)

	return nil
}

func (c *Channel) Dequeue() (*task.Task, error) {
	if c.chanel == nil {
		return nil, errors.New("channel not created yet")
	}

	t := <-c.chanel

	c.wg.Done()

	fmt.Printf("Task with ID %d and name %s finished the processing\n", t.ID, t.Name)

	fmt.Printf("Queue size: %d\n", c.Size())
	fmt.Printf("The queue got empty?: %t\n", c.IsEmpty())

	return t, nil
}

func (c *Channel) IsEmpty() bool {
	return len(c.chanel) == 0
}

func (c *Channel) Size() int {
	return len(c.chanel)
}
