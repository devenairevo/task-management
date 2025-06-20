package main

type TaskManager interface {
	CreateTask(*CreateTaskParams) (*Task, error) // Creates a new task
	ListTasks() ([]*Task, error)                 // Lists all tasks
	DescribeTask(taskID int) (*Task, error)      // Provides details for a specific task
	UpdateTask(*Task) error                      // Updates an existing task
}

type CreateTaskParams struct {
	UserID   int
	UserName string
	TaskID   int
	Name     string
}
