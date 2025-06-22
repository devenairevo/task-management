package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/devenairevo/task-management/internal/contracts/queuer"
	"github.com/devenairevo/task-management/internal/contracts/tasker"
	"github.com/devenairevo/task-management/test/data"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func PostTasksHandler(taskMng tasker.Tasker, channel queuer.Queuer, wg *sync.WaitGroup) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Println("Adding tasks to the queue....")

			tasksList := data.GenerateMockTasks(taskMng)

			for _, t := range tasksList {
				err := channel.Enqueue(t)
				if err != nil {
					return
				}
			}

			wg.Wait()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			if err := json.NewEncoder(w).Encode(tasksList); err != nil {
				http.Error(w, "Issue with encoding", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		return
	}
}

func GetTasksHandler(taskMng tasker.Tasker, channel queuer.Queuer) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Println("Getting all tasks....")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			tasks, err := taskMng.ListTasks()
			if err != nil {
				http.Error(w, "Issue with the task lists", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

			if err := json.NewEncoder(w).Encode(tasks); err != nil {
				http.Error(w, "Issue with encoding", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}
			for _, v := range tasks {
				fmt.Printf("Task with an id %d, name '%s' and status - %s\n", v.ID, v.Name, v.Status)
			}

			fmt.Printf("Queue size: %d\n", channel.Size())
			fmt.Printf("The queue got empty?: %t\n", channel.IsEmpty())
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		return
	}

}

func TaskByID(taskMng tasker.Tasker, channel queuer.Queuer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 3 {
			http.Error(w, "Incorrect path", http.StatusBadRequest)
			return
		}

		id := pathParts[2]
		taskID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}
		fmt.Println(taskID)
		switch r.Method {
		case http.MethodGet:
			fmt.Println("Getting task status....")
			task, err := taskMng.DescribeTask(taskID)
			if err != nil {
				http.Error(w, "Issue with the getting task", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(task)
			if err != nil {
				http.Error(w, "Issue with the encoding task", http.StatusInternalServerError)
				fmt.Println(err)
			}

			fmt.Printf("Found ✅ : %s with the ID %d and status - %s\n", task.Name, task.ID, task.Status)
		case http.MethodPut:
			task, err := taskMng.DescribeTask(taskID)

			fmt.Println("Updating task....")

			err = taskMng.UpdateTask(task)
			if err != nil {
				http.Error(w, "Issue with the updating task", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(task)
			if err != nil {
				http.Error(w, "Issue with the encoding task", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

			fmt.Printf("Queue size: %d\n", channel.Size())
			fmt.Printf("The queue got empty?: %t\n", channel.IsEmpty())

			fmt.Printf("Found ✅ : %s with the ID %d and status - %s\n", task.Name, task.ID, task.Status)
		}

		return
	}
}
