package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/MarySmirnova/tasks/internal/config"
	"github.com/MarySmirnova/tasks/pkg/storage/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ctx context.Context = context.Background()

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(cfg config.Postgres) (*Storage, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := pgxpool.Connect(ctx, connString)
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

func (s *Storage) NewTask(task models.Task) (int, error) {
	query := `
		INSERT INTO tasks (
			opened,
			author_id,
			assigned_id,
			title,
			content) 
		VALUES ($1, $2, $3, $4, $5)	
		RETURNING id;`

	openTime := time.Now()

	var id int
	err := s.db.QueryRow(ctx, query, openTime, task.Author.ID, task.Assigned.ID, task.Title, task.Content).Scan(&id)
	if err != nil {
		return 0, err
	}

	query = `INSERT INTO tasks_labels 
		(task_id, label_id)
		VALUES ($1, $2);`

	for _, label := range task.Label {
		_, err := s.db.Exec(ctx, query, id, label.ID)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

//Для поиска всех записей нужно передавать новую пустую структуру, где User.ID будет равен 0.
func (s *Storage) GetTasks(author models.User) ([]models.Task, error) {
	query := `
		SELECT
			id,
    		opened,
    		closed,
    		author_id,
    		assigned_id,
    		title,
    		content 
		FROM tasks
		WHERE
			($1 = 0 OR author_id = $1);`

	rows, err := s.db.Query(ctx, query, author.ID)
	if err != nil {
		return nil, err
	}

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		err = rows.Scan(&t.ID, &t.Opened, &t.Closed, &t.AuthorID, &t.AssignedID, &t.Title, &t.Content)
		if err != nil {
			return nil, err
		}

		queryUser := `SELECT name FROM users WHERE id = &1;`
		if t.AuthorID != 0 {
			t.Author = models.User{}
			t.Author.ID = t.AuthorID
			err = s.db.QueryRow(ctx, queryUser, t.AuthorID).Scan(&t.Author.Name)
			if err != nil {
				return nil, err
			}
		}

		if t.AssignedID != 0 {
			t.Assigned = models.User{}
			t.Assigned.ID = t.AssignedID
			err = s.db.QueryRow(ctx, queryUser, t.AssignedID).Scan(&t.Assigned.Name)
			if err != nil {
				return nil, err
			}
		}

		queryLabel := `
		SELECT tasks_labels.label_id, labels.name
		FROM tasks_labels, labels
		WHERE tasks_labels.label_id = labels.id
		AND tasks_labels.task_id = $1;`

		labels, err := s.db.Query(ctx, queryLabel, t.ID)
		if err != nil {
			return nil, err
		}

		t.Label = []models.Label{}
		for labels.Next() {
			var l models.Label
			err = labels.Scan(&l.ID, &l.Name)
			if err != nil {
				return nil, err
			}
		}
		if err = labels.Err(); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (s *Storage) UpdateTask(id int) error {
	return nil
}

func (s *Storage) DeleteTask(id int) error {
	return nil
}
