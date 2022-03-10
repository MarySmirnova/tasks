package models

import (
	"github.com/jackc/pgtype"
)

type Task struct {
	ID       int              `json:"id"`
	Opened   pgtype.Timestamp `json:"opened"`
	Closed   pgtype.Timestamp `json:"closed"`
	Author   User             `json:"author"`
	Assigned User             `json:"assigned"`
	Title    string           `json:"title"`
	Content  string           `json:"content"`
	Label    []Label          `json:"label"`

	AuthorID   int
	AssignedID int
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Label struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
