package storage

import "github.com/MarySmirnova/tasks/pkg/storage/models"

type InterfaceDB interface {
	NewTask(models.Task) (int, error)
	GetAllTasks() ([]models.Task, error)
	GetTasks(author models.User, label models.Label) ([]models.Task, error)
	UpdateTask(id int) error
	DeleteTask(id int) error
}
