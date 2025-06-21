package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

/*
Develop a task management mechanism accessible via HTTP.
The system should provide the following capabilities:

- Users can create new tasks via an HTTP endpoint.
- Newly created tasks should be added to a queue for asynchronous processing.
- Users can list all tasks and check the status of individual tasks using their task ID.
- The system must include both a task queue and a task management component to handle task execution and status tracking.

*/

func main() {
	server := &http.Server{
		Addr:         ":2025",
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}

	wg := &sync.WaitGroup{}
	channel, _ := NewQueueManager("local", wg)
	taskMng, _ := NewTaskManager("local")

	http.HandleFunc("/tasks", getTasksHandler(taskMng, channel))
	http.HandleFunc("/task/create", postTasksHandler(taskMng, channel, wg))
	http.HandleFunc("/tasks/", taskByID(taskMng, channel))

	// Run goroutine for Dequeue
	go func(channel QueueManager, taskMng TaskManager) {
		for {
			finishedTask, err := channel.Dequeue()
			if err != nil {
				continue
			}

			task, err := taskMng.DescribeTask(finishedTask.ID)
			if err == nil {
				task.Status = Done
			}
		}
	}(channel, taskMng)

	fmt.Printf("Server started, please make your HTTP requests to the localhost with a port %s and watch the results in terminal....\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func NewQueueManager(queueDriver string, wg *sync.WaitGroup) (QueueManager, error) {
	switch queueDriver {
	case "local":
		return NewChannel(buffSize, wg), nil
	default:
		return nil, fmt.Errorf("unsupported queue driver: %s", queueDriver)
	}
}

func NewTaskManager(taskDriver string) (TaskManager, error) {
	switch taskDriver {
	case "local":
		return NewLocalTaskManager(), nil
	default:
		return nil, fmt.Errorf("unsupported task driver: %s", taskDriver)
	}
}
