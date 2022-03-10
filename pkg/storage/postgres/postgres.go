package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/MarySmirnova/tasks/internal/config"
	"github.com/MarySmirnova/tasks/pkg/storage/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ctx context.Context = context.Background()
var nowTime = time.Now()

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

func (s *Storage) GetPGPool() *pgxpool.Pool {
	return s.db
}

type StorTX struct {
	tx pgx.Tx
}

func (s *Storage) NewStorTX() (*StorTX, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &StorTX{
		tx: tx,
	}, nil
}

//NewTask - создает новую задачу
func (s *Storage) NewTask(task models.Task) (int, error) {
	tx, err := s.NewStorTX()
	if err != nil {
		return 0, err
	}
	defer tx.tx.Rollback(ctx)

	query := `
		INSERT INTO tasks (
			opened,
			author_id,
			assigned_id,
			title,
			content) 
		VALUES ($1, $2, $3, $4, $5)	
		RETURNING id;`

	var id int
	err = tx.tx.QueryRow(ctx, query, nowTime, task.Author.ID, task.Assigned.ID, task.Title, task.Content).Scan(&id)
	if err != nil {
		return 0, err
	}

	query = `INSERT INTO tasks_labels 
		(task_id, label_id)
		VALUES ($1, $2);`

	for _, label := range task.Label {
		_, err := tx.tx.Exec(ctx, query, id, label.ID)
		if err != nil {
			return 0, err
		}
	}

	tx.tx.Commit(ctx)
	return id, nil
}

//GetTasks - возвращает выборку задач по юзеру.
//Для поиска всех записей нужно передавать новую пустую структуру, где User.ID будет равен 0.
func (s *Storage) GetTasks(author models.User) ([]*models.Task, error) {
	tx, err := s.NewStorTX()
	if err != nil {
		return nil, err
	}
	defer tx.tx.Rollback(ctx)

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

	rows, err := tx.tx.Query(ctx, query, author.ID)
	if err != nil {
		return nil, err
	}

	tasks := []*models.Task{}
	for rows.Next() {
		var t models.Task
		err = rows.Scan(&t.ID, &t.Opened, &t.Closed, &t.AuthorID, &t.AssignedID, &t.Title, &t.Content)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	err = tx.fillTasks(tasks)
	if err != nil {
		return nil, err
	}

	tx.tx.Commit(ctx)
	return tasks, nil
}

//GetTasksByLabel - возвращает выборку задач по метке.
func (s *Storage) GetTasksByLabel(label models.Label) ([]*models.Task, error) {
	tx, err := s.NewStorTX()
	if err != nil {
		return nil, err
	}
	defer tx.tx.Rollback(ctx)

	query := `
	SELECT task_id
	FROM tasks_labels
	WHERE label_id = $1;`

	rows, err := tx.tx.Query(ctx, query, label.ID)
	if err != nil {
		return nil, err
	}

	tasksID := []int{}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		tasksID = append(tasksID, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	query = `
	SELECT
		id,
		opened,
		closed,
		author_id,
		assigned_id,
		title,
		content 
	FROM tasks
	WHERE id = $1;`

	var tasks []*models.Task
	for _, id := range tasksID {
		var t models.Task
		err = tx.tx.QueryRow(ctx, query, id).Scan(&t.ID, &t.Opened, &t.Closed, &t.AuthorID, &t.AssignedID, &t.Title, &t.Content)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	err = tx.fillTasks(tasks)
	if err != nil {
		return nil, err
	}

	tx.tx.Commit(ctx)
	return tasks, nil
}

//UpdateTask - обновляет данные в задаче
//(но не устанавливает дату закрытия и не обновляет метки, для этого нужно будет писать отдельные методы)
func (s *Storage) UpdateTask(task models.Task) error {
	query := `
	UPDATE tasks SET 
		author_id = $1, 
		assigned_id = $2,
		title = $3,
		content = $4
	WHERE id = $5;`

	_, err := s.db.Exec(ctx, query, task.Author.ID, task.Assigned.ID, task.Title, task.Content, task.ID)
	if err != nil {
		return err
	}

	return nil
}

//DeleteTask - удаляет задачу
func (s *Storage) DeleteTask(id int) error {
	tx, err := s.NewStorTX()
	if err != nil {
		return err
	}
	defer tx.tx.Rollback(ctx)

	query := `
	DELETE FROM tasks_labels
	WHERE task_id = $1`

	_, err = tx.tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	query = `
	DELETE FROM tasks
	WHERE id = $1`

	_, err = tx.tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	tx.tx.Commit(ctx)
	return nil
}

func (tx *StorTX) fillTasks(tasks []*models.Task) error {
	for _, t := range tasks {
		queryUser := `SELECT name FROM users WHERE id = $1;`
		if t.AuthorID != 0 {
			t.Author = models.User{}
			t.Author.ID = t.AuthorID
			err := tx.tx.QueryRow(ctx, queryUser, t.AuthorID).Scan(&t.Author.Name)
			if err != nil {
				return err
			}
		}

		if t.AssignedID != 0 {
			t.Assigned = models.User{}
			t.Assigned.ID = t.AssignedID
			err := tx.tx.QueryRow(ctx, queryUser, t.AssignedID).Scan(&t.Assigned.Name)
			if err != nil {
				return err
			}
		}

		queryLabel := `
		SELECT tasks_labels.label_id, labels.name
		FROM tasks_labels
		JOIN labels ON (tasks_labels.label_id = labels.id)
		WHERE tasks_labels.task_id = $1;`

		labels, err := tx.tx.Query(ctx, queryLabel, t.ID)
		if err != nil {
			return err
		}

		t.Label = []models.Label{}
		for labels.Next() {
			var l models.Label
			err = labels.Scan(&l.ID, &l.Name)
			if err != nil {
				return err
			}
			t.Label = append(t.Label, l)
		}
		if err = labels.Err(); err != nil {
			return err
		}
	}
	return nil
}
