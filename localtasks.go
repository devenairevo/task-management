package main

import (
	"errors"
	"fmt"
)

type Status string

const (
	Created    Status = "Created"
	Processing Status = "Processing"
	Updated    Status = "Updated"
	Done       Status = "Done"
)

type Task struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status"`
}

func NewTask(id int, name string, status Status) (*Task, error) {
	if id != 0 && name != "" {
		return &Task{id, name, status}, nil
	}

	return nil, errors.New("something wrong with creating a new task")

}

type UserTask struct {
	*User
	*Task
}

func (lt *LocalTaskManager) CreateTask(params *CreateTaskParams) (*Task, error) {
	if params.UserID <= 0 || params.TaskID == 0 || params.Name == "" {
		return nil, errors.New("couldn't create user with task")
	}

	task, err := NewTask(params.TaskID, params.Name, Created)
	if err != nil {
		return nil, err
	}

	user, err := NewUser(params.UserID, params.UserName)
	if err != nil {
		return nil, err
	}

	userTask, err := NewUserTask(user, task)
	if err != nil {
		return nil, err
	}

	lt.AddTask(userTask)

	return task, nil
}

func NewUserTask(user *User, task *Task) (*UserTask, error) {
	if user.ID > 0 && task.ID != 0 && task.Name != "" {
		return &UserTask{
			User: user,
			Task: task,
		}, nil
	}
	return nil, errors.New("couldn't create user with task")
}

type LocalTaskManager struct {
	UserTasks []*UserTask
}

func NewLocalTaskManager() *LocalTaskManager {
	return &LocalTaskManager{
		UserTasks: []*UserTask{},
	}
}

func (lt *LocalTaskManager) AddTask(userTasks ...*UserTask) {
	for _, v := range userTasks {
		lt.UserTasks = append(lt.UserTasks, v)
	}
}

func (lt *LocalTaskManager) GetUserTasks(userID int) []*Task {
	var userTasks []*Task

	for _, userTask := range lt.UserTasks {
		if userTask.User.ID == userID {
			userTasks = append(userTasks, userTask.Task)
		}
	}

	return userTasks
}

func (lt *LocalTaskManager) DescribeTask(taskID int) (*Task, error) {
	for _, userTask := range lt.UserTasks {
		if userTask.Task.ID == taskID {
			return userTask.Task, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("❌  Task with ID %d not found", taskID))
}

func (lt *LocalTaskManager) ListTasks() ([]*Task, error) {
	if len(lt.UserTasks) == 0 {
		return nil, errors.New("no tasks found")
	}

	var tasks []*Task
	for _, userTask := range lt.UserTasks {
		if userTask.Task != nil {
			tasks = append(tasks, userTask.Task)
		}
	}

	return tasks, nil
}

func (lt *LocalTaskManager) UpdateTask(task *Task) error {
	if task == nil || task.ID < 0 {
		return errors.New("issue for the updating current task")
	}
	task.Status = Updated

	return nil
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
