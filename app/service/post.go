package service

import (
	"context"
	"fmt"
	"time"

	"stepik.leoscode.http/app/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePost(userID string, threadID int, cont string) (*models.Post, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return nil, fmt.Errorf("Connection error %s", err)
	}
	var username string
	if err := conn.QueryRow(ctx, `SELECT username FROM users WHERE id = $1`, userID).Scan(&username); err != nil {
		return nil, fmt.Errorf("User not found")
	}
	var isLocked bool
	if err := conn.QueryRow(ctx, `SELECT is_locked FROM threads WHERE id = $1`, threadID).Scan(&isLocked); err != nil {
		return nil, fmt.Errorf("User not found")
	}
	if isLocked {
		return nil, fmt.Errorf("Thread is locked")
	}
	post := &models.Post{
		AuthorID:  userID,
		ThreadID:  threadID,
		Content:   cont,
		CreatedAt: time.Now(),
	}
	var postID int
	err = conn.QueryRow(ctx, `INSERT INTO posts (author_id, thread_id, content)
		VALUES ($1, $2, $3)
		RETURNING id`, post.AuthorID, post.ThreadID, post.Content).Scan(&postID)
	if err != nil {
		return nil, err
	}
	post.Id = postID

	return post, nil
}

func GetPosts() {

}
