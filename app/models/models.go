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

	IsLocked  bool       `json:"is_locked"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Parameters struct {
	Limit  int
	Offset int
	Tag    string
	Author string
	Sort   string
}

type Contents struct {
	Content string `json:"content"`
}

type Post struct {
	Id       int
	AuthorID string
	ThreadID int
	Content  string `json:"content"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Attachments struct {
	Id        string
	ThreadId  int
	Filename  string
	MimeType  string
	Size      int
	Path      string `json:"-"`
	CreatedAt time.Time
}
