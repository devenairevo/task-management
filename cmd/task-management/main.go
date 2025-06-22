package main

import (
	"fmt"
	"github.com/devenairevo/task-management/internal/app"
	"github.com/devenairevo/task-management/internal/handlers"
	"github.com/devenairevo/task-management/internal/queue"
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

const buffSize = 5

func main() {
	server := &http.Server{
		Addr:         ":2025",
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}

	wg := &sync.WaitGroup{}
	channel, _ := app.NewQueueManager(buffSize, "local", wg)
	taskMng, _ := app.NewTaskManager("local")

	http.HandleFunc("/tasks", handlers.GetTasksHandler(taskMng, channel))
	http.HandleFunc("/task/create", handlers.PostTasksHandler(taskMng, channel, wg))
	http.HandleFunc("/tasks/", handlers.TaskByID(taskMng, channel))

	// Run goroutine for Dequeue
	queue.StartQueueWorker(channel, taskMng)
	
	fmt.Printf("Server started, please make your HTTP requests "+
		"to the localhost with a port %s and watch the results in terminal....\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
