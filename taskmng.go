package main

import (
	"errors"
	"fmt"
)

type Status string

const (
	Created    Status = "Created"
	Processing Status = "Processing"
	Done       Status = "Done"
)

type Task struct {
	ID     int
	Name   string
	Status Status
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

func NewUserTask(user *User, task *Task) (*UserTask, error) {
	if user.ID > 0 && task.ID != 0 && task.Name != "" {
		return &UserTask{
			User: user,
			Task: task,
		}, nil
	}
	return nil, errors.New("couldn't create user with task")
}

type TaskManager struct {
	UserTasks []*UserTask
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		UserTasks: []*UserTask{},
	}
}

func (tm *TaskManager) AddTask(userTasks ...*UserTask) {
	for _, v := range userTasks {
		tm.UserTasks = append(tm.UserTasks, v)
	}
}

func (tm *TaskManager) GetUserTasks(userID int) []*Task {
	var userTasks []*Task

	for _, userTask := range tm.UserTasks {
		if userTask.User.ID == userID {
			userTasks = append(userTasks, userTask.Task)
		}
	}

	return userTasks
}

func (tm *TaskManager) GetTaskByID(taskID int) (*Task, error) {
	for _, userTask := range tm.UserTasks {
		if userTask.Task.ID == taskID {
			return userTask.Task, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("❌  Task with ID %d not found", taskID))
}
