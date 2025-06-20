package main

type QueueManager interface {
	Enqueue(*Task) error
	Dequeue() (*Task, error)
	IsEmpty() bool
	Size() int
}
