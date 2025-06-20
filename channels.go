package main

import (
	"errors"
	"fmt"
	"sync"
)

const buffSize = 5

type Channel struct {
	mu       sync.Mutex
	wg       *sync.WaitGroup
	buffSize int
	Task     *Task
	Chan     chan *Task
}

func NewChannel(buffSize int, wg *sync.WaitGroup) *Channel {
	return &Channel{
		buffSize: buffSize,
		Chan:     make(chan *Task, buffSize),
		wg:       wg,
	}
}

func (c *Channel) Enqueue(task *Task) error {
	if task.ID <= 0 || task.Name == "" {
		return errors.New("error with adding task to the queue")
	}
	if c.Chan == nil {
		return errors.New("channel not created yet")
	}

	c.wg.Add(1)
	c.Chan <- task

	task.Status = Created
	fmt.Printf("Task with ID %d and name %s created\n", task.ID, task.Name)

	task.Status = Processing
	fmt.Printf("Task with ID %d and name %s started processing\n", task.ID, task.Name)

	return nil
}

func (c *Channel) Dequeue() (*Task, error) {
	if c.Chan == nil {
		return nil, errors.New("channel not created yet")
	}

	task := <-c.Chan

	c.wg.Done()

	fmt.Printf("Task with ID %d and name %s finished the processing\n", task.ID, task.Name)

	fmt.Printf("Queue size: %d\n", c.Size())
	fmt.Printf("The queue got empty?: %t\n", c.IsEmpty())

	return task, nil
}

func (c *Channel) IsEmpty() bool {
	return len(c.Chan) == 0
}

func (c *Channel) Size() int {
	return len(c.Chan)
}
