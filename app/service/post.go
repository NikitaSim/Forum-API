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
	defer conn.Close()
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

func GetPosts(id int, param models.Parameters) ([]models.Post, int, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return nil, 0, fmt.Errorf("Connection error %s", err)
	}
	defer conn.Close()
	var total int
	err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM posts`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	posts := make([]models.Post, 0)
	rows, err := conn.Query(ctx,
		`SELECT id, author_id, thread_id, content, created_at, updated_at FROM posts
		WHERE thread_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`, id, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.AuthorID, &post.ThreadID, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, 0, err
		}
		posts = append(posts, post)
	}
	rows.Close()

	return posts, total, nil
}

func DeletePost(userID string, id int) error {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return fmt.Errorf("Connection error %s", err)
	}
	defer conn.Close()

	var postOwner string
	if err := conn.QueryRow(ctx, `SELECT author_id FROM posts WHERE id = $1`, id).Scan(&postOwner); err != nil {
		return err
	}
	if postOwner != userID {
		return fmt.Errorf("This user can not delete this post")
	}

	if _, err := conn.Exec(ctx, `DELETE FROM posts WHERE id = $1`, id); err != nil {
		return err
	}

	return nil
}

func UpdatePost(userID string, id int, content string) (*time.Time, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return nil, fmt.Errorf("Connection error %s", err)
	}
	defer conn.Close()

	var postOwner string
	var isLocked bool
	if err := conn.QueryRow(ctx, `SELECT posts.author_id, threads.is_locked FROM posts
	INNER JOIN threads ON posts.thread_id = threads.id WHERE posts.id = $1`, id).Scan(&postOwner, &isLocked); err != nil {
		return nil, err
	}
	if postOwner != userID || isLocked {
		return nil, fmt.Errorf("This user can not update this post or thread is locked")
	}

	var updatedAt *time.Time
	if err := conn.QueryRow(ctx, `UPDATE posts SET
	content = $1,
	updated_at = NOW()
	WHERE id = $2
	RETURNING updated_at`, content, id).Scan(&updatedAt); err != nil {
		return nil, err
	}

	return updatedAt, nil
}
