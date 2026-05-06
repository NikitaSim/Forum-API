package models

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
