package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	queue := NewQueue(5)
	taskMng := NewTaskManager()

	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Println("Adding tasks to the queue....")
			generateMockTasks(taskMng)
			for _, task := range taskMng.UserTasks {
				queue.Chan <- task
			}
		}

		return
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Println("Getting all users tasks....")
			for _, t := range taskMng.UserTasks {
				userTasks := taskMng.GetUserTasks(t.User.ID)

				for _, task := range userTasks {
					fmt.Printf("The user with id %d has task name '%s' and status - %s\n", t.User.ID, task.Name, task.Status)
				}
			}

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
			task, err := taskMng.GetTaskByID(taskID)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Found ✅ : %s with the ID %d and status - %s\n", task.Name, task.ID, task.Status)
		}

		return
	})

	go processQueue(queue.Chan)
	fmt.Printf("Server started, please make your HTTP requests to the localhost with a port %s and watch the results in terminal....\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func processQueue(queue chan *UserTask) {
	for value := range queue {
		value.Task.Status = Processing
		fmt.Printf("Task %d for %s started processing\n", value.Task.ID, value.User.Name)
		time.Sleep(3 * time.Second)
		value.Task.Status = Done
		fmt.Printf("Task %d for %s finished the processing\n", value.Task.ID, value.User.Name)

	}
}

func generateMockTasks(manager *TaskManager) {
	user1, _ := NewUser(1, "User1")
	task1, _ := NewTask(1, "Deploy service", Created)
	task11, _ := NewTask(5, "Deploy another service", Created)

	user2, _ := NewUser(2, "User2")
	task2, _ := NewTask(2, "Run integration tests", Created)

	user3, _ := NewUser(3, "User3")
	task3, _ := NewTask(3, "Generate report", Created)

	ut1, _ := NewUserTask(user1, task1)
	ut2, _ := NewUserTask(user1, task11)
	ut3, _ := NewUserTask(user2, task2)
	ut4, _ := NewUserTask(user3, task3)

	manager.AddTask(ut1, ut2, ut3, ut4)
}
