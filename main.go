package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Println("Adding tasks to the queue....")

			tasksList := generateMockTasks(taskMng)

			for _, task := range tasksList {
				err := channel.Enqueue(task)
				if err != nil {
					return
				}
			}

			wg.Wait()

		}

		return
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Println("Getting all tasks....")

			tasks, err := taskMng.ListTasks()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, v := range tasks {

				fmt.Printf("Task with an id %d, name '%s' and status - %s\n", v.ID, v.Name, v.Status)

			}

			fmt.Printf("Queue size: %d\n", channel.Size())
			fmt.Printf("The queue got empty?: %t\n", channel.IsEmpty())

		}

		return
	})

	http.HandleFunc("/task/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			path := r.URL.Path

			id := path[len("/task/"):]
			taskID, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("Getting task status....")
			task, err := taskMng.DescribeTask(taskID)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Found ✅ : %s with the ID %d and status - %s\n", task.Name, task.ID, task.Status)
		}

		if r.Method == http.MethodPut {
			path := r.URL.Path

			id := path[len("/task/"):]
			taskID, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println(err)
			}
			task, err := taskMng.DescribeTask(taskID)

			fmt.Println("Updating task....")

			err = taskMng.UpdateTask(task)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Queue size: %d\n", channel.Size())
			fmt.Printf("The queue got empty?: %t\n", channel.IsEmpty())

			fmt.Printf("Found ✅ : %s with the ID %d and status - %s\n", task.Name, task.ID, task.Status)
		}

		return
	})

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

func generateMockTasks(manager TaskManager) []*Task {
	var tasksList []*Task

	task1 := &CreateTaskParams{1, "User1", 1, "Deploy service"}
	ut1, _ := manager.CreateTask(task1)

	task2 := &CreateTaskParams{1, "User1", 2, "Deploy another service"}
	ut2, _ := manager.CreateTask(task2)

	task3 := &CreateTaskParams{2, "User2", 3, "Run integration tests"}
	ut3, _ := manager.CreateTask(task3)

	task4 := &CreateTaskParams{3, "User3", 4, "Generate report"}
	ut4, _ := manager.CreateTask(task4)

	tasksList = append(tasksList, ut1, ut2, ut3, ut4)

	return tasksList
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
