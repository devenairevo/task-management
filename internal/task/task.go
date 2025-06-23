package task

import (
	"errors"
	"fmt"
	"github.com/devenairevo/task-management/internal/types"
	"github.com/devenairevo/task-management/internal/user"
)

type Task struct {
	ID     int          `json:"id"`
	Name   string       `json:"name"`
	Status types.Status `json:"status"`
}

type CreateTaskParams struct {
	UserID   int
	UserName string
	TaskID   int
	Name     string
}

func NewTask(id int, name string, status types.Status) (*Task, error) {
	if id != 0 && name != "" {
		return &Task{ID: id, Name: name, Status: status}, nil
	}

	return nil, errors.New("something wrong with creating a new task")

}

func (lt *LocalTaskManager) CreateTask(params *CreateTaskParams) (*Task, error) {
	if params.UserID <= 0 || params.TaskID == 0 || params.Name == "" {
		return nil, errors.New("couldn't create user with task")
	}

	task, err := NewTask(params.TaskID, params.Name, types.Created)
	if err != nil {
		return nil, err
	}

	u, err := user.New(params.UserID, params.UserName)
	if err != nil {
		return nil, err
	}

	userTask, err := NewUserTask(u, task)
	if err != nil {
		return nil, err
	}

	lt.AddTask(userTask)

	return task, nil
}

func NewUserTask(user *user.User, task *Task) (*UserTask, error) {
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

	return nil, fmt.Errorf("❌  Task with ID %d not found", taskID)
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
		return errors.New("no task found")
	}
	task.Status = types.Updated

	return nil
}
