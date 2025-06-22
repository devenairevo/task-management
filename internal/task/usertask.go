package task

import (
	"github.com/devenairevo/task-management/internal/user"
)

type UserTask struct {
	*user.User
	*Task
}
