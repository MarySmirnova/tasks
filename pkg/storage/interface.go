package storage

import "github.com/MarySmirnova/tasks/pkg/storage/models"

type InterfaceDB interface {
	NewTask(models.Task) (int, error)
	GetTasks(author models.User) ([]*models.Task, error)
	GetTasksByLabel(label models.Label) ([]*models.Task, error)
	UpdateTask(task models.Task) error
	DeleteTask(id int) error
}
