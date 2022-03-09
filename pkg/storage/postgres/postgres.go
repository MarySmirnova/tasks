package postgres

import (
	"context"

	"github.com/MarySmirnova/tasks/pkg/storage/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ctx context.Context = context.Background()

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(cfg string) (*Storage, error) {

	db, err := pgxpool.Connect(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(ctx); err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) NewTask(models.Task) error {
	return nil
}

func (s *Storage) GetAllTasks() ([]models.Task, error) {
	return nil, nil
}

func (s *Storage) GetTasks(author models.User, label models.Label) ([]models.Task, error) {
	return nil, nil
}

func (s *Storage) UpdateTask(id int) error {
	return nil
}

func (s *Storage) DeleteTask(id int) error {
	return nil
}
