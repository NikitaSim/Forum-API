package models

import "time"

const PostgresqlConnString = "postgres://user:password@localhost:5433/postgres"

var DBTables = []string{
	"users",
	"threads",
	"thread_tags",
	"posts",
	"attachments",
	"idempotency_keys",
}

type Users struct {
	Name     string `form:"username" binding:"required,min=3,max=32"`
	Password string `form:"password" binding:"required,min=8,max=64"`
}

type Thread struct {
	Id       int
	AuthorID string
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`

	IsLocked  bool
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Parameters struct {
	Limit  int
	Offset int
	Tag    string
	Author string
}
