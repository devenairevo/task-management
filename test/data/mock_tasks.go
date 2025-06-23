package data

import (
	"github.com/devenairevo/task-management/internal/contracts/tasker"
	"github.com/devenairevo/task-management/internal/task"
)

func GenerateMockTasks(manager tasker.Tasker) []*task.Task {
	var tasksList []*task.Task

	task1 := &task.CreateTaskParams{UserID: 1, UserName: "User1", TaskID: 1, Name: "Deploy service"}
	ut1, _ := manager.CreateTask(task1)

	task2 := &task.CreateTaskParams{UserID: 1, UserName: "User1", TaskID: 2, Name: "Deploy another service"}
	ut2, _ := manager.CreateTask(task2)

	task3 := &task.CreateTaskParams{UserID: 2, UserName: "User2", TaskID: 3, Name: "Run integration tests"}
	ut3, _ := manager.CreateTask(task3)

	task4 := &task.CreateTaskParams{UserID: 3, UserName: "User3", TaskID: 4, Name: "Generate report"}
	ut4, _ := manager.CreateTask(task4)

	tasksList = append(tasksList, ut1, ut2, ut3, ut4)

	return tasksList
}
