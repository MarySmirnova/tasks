package models

import (
	"github.com/jackc/pgtype"
)

type Task struct {
	//json:
	ID       int
	Opened   pgtype.Timestamp
	Closed   pgtype.Timestamp
	Author   User
	Assigned User
	Title    string
	Content  string
	Label    []Label

	//not json:
	AuthorID   int
	AssignedID int
}

type User struct {
	ID   int
	Name string
}

type Label struct {
	ID   int
	Name string
}
