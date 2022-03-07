package models

import "time"

type Task struct {
	ID       int
	Opened   time.Time
	Closed   time.Time
	Author   User
	Assigned User
	Title    string
	Content  string
	Label    []Label
}

type User struct {
	ID   int
	Name string
}

type Label struct {
	ID   int
	Name string
}
