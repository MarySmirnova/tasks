package memdb

import "github.com/MarySmirnova/tasks/pkg/storage/models"

type DB []models.Task

func (db DB) NewTask(models.Task) (int, error) {
	return 0, nil
}

func (db DB) GetTasks(author models.User) ([]*models.Task, error) {
	return nil, nil
}

func (db DB) GetTasksByLabel(label models.Label) ([]*models.Task, error) {
	return nil, nil
}

func (db DB) UpdateTask(task models.Task) error {
	return nil
}

func (db DB) DeleteTask(id int) error {
	return nil
}
