package main

type Queue struct {
	buffSize int
	UserTask *UserTask
	Chan     chan *UserTask
}

func NewQueue(buffSize int) *Queue {
	return &Queue{
		buffSize: buffSize,
		Chan:     make(chan *UserTask, buffSize),
	}
}
