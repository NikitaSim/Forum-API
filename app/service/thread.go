package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"stepik.leoscode.http/app/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func hashThread(thread models.Thread) (string, error) {
	data, err := json.Marshal(thread)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}

func CreateThread(userID, idempotencyKey string, thread models.Thread) (models.Thread, int, error) {

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, models.PostgresqlConnString)
	if err != nil {
		fmt.Println("Connection error", err)
		return models.Thread{}, 500, err
	}
	defer conn.Close(ctx)

	hash, err := hashThread(thread)
	if err != nil {
		return models.Thread{}, 500, err
	}

	var hashItem string
	var responseBody []byte
	var statusCode int

	err = conn.QueryRow(ctx, "SELECT request_hash, response_body, status_code FROM idempotency_keys WHERE user_id = $1 and key = $2", userID, idempotencyKey).Scan(&hashItem, &responseBody, &statusCode)
	if err == nil {

		if hashItem != hash {
			return models.Thread{}, 409, fmt.Errorf("conflict")
		}

		var oldThread models.Thread
		err := json.Unmarshal(responseBody, &oldThread)
		if err != nil {
			return models.Thread{}, 500, err
		}
		return oldThread, 200, nil
	}

	var threadId int

	err = conn.QueryRow(ctx, "INSERT INTO threads (author_id, title, content) VALUES ($1, $2, $3)  RETURNING id", userID, thread.Title, thread.Content).Scan(&threadId)
	if err != nil {
		return models.Thread{}, 500, fmt.Errorf("insert error: %s", err)
	}

	for _, tag := range thread.Tags {
		_, err = conn.Exec(ctx, `INSERT INTO thread_tags (thread_id, tag) VALUES ($1, $2)`, threadId, tag)
		if err != nil {
			return models.Thread{}, 500, err
		}
	}

	thread.AuthorID = userID
	thread.Id = threadId
	thread.CreatedAt = time.Now()
	respByte, err := json.Marshal(thread)
	if err != nil {
		return models.Thread{}, 500, err
	}

	_, err = conn.Exec(ctx, `INSERT INTO idempotency_keys
	(user_id, key, request_hash, response_body, status_code) VALUES
	($1, $2, $3, $4, $5)`, userID, idempotencyKey, hash, respByte, 201)
	if err != nil {
		return models.Thread{}, 500, err
	}

	return thread, 201, nil
}

func GetThreads(param models.Parameters) ([]models.Thread, int, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString) //.Connect(ctx, models.PostgresqlConnString)
	if err != nil {
		fmt.Println("Connection error", err)
		return nil, 0, err
	}
	defer conn.Close()
	var total int
	err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM threads`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	threads := make([]models.Thread, 0)
	if param.Author == "" && param.Tag == "" {
		query := `SELECT threads.id,
			threads.author_id,
			threads.title,
			threads.content,
			threads.is_locked,
			threads.created_at,
			threads.updated_at,
			array_agg(thread_tags.tag) as tags
			FROM threads LEFT JOIN thread_tags ON threads.id = thread_tags.thread_id
			GROUP BY threads.id 
			ORDER BY created_at `
		if param.Sort == "old" {
			query += "ASC LIMIT $1 OFFSET $2"
		} else {
			query += "DESC LIMIT $1 OFFSET $2"
		}

		rows, err := conn.Query(ctx, query, param.Limit, param.Offset)
		if err != nil {
			return nil, 0, err
		}

		for rows.Next() {
			tags := make([]string, 0)
			var thread models.Thread
			if err := rows.Scan(&thread.Id, &thread.AuthorID, &thread.Title, &thread.Content, &thread.IsLocked, &thread.CreatedAt, &thread.UpdatedAt, &tags); err != nil {
				return nil, 0, err
			}
			thread.Tags = tags
			threads = append(threads, thread)
		}

		return threads, total, nil
	}

	if param.Author != "" {
		query := `SELECT
			threads.id,
			threads.author_id,
			threads.title,
			threads.content,
			threads.is_locked,
			threads.created_at,
			threads.updated_at,
			array_agg(thread_tags.tag) as tags
		FROM threads
		LEFT JOIN thread_tags
		ON threads.id = thread_tags.thread_id
		WHERE threads.author_id = $1
		GROUP BY
			threads.id
		ORDER BY created_at `
		if param.Sort == "old" {
			query += "ASC LIMIT $2 OFFSET $3"
		} else {
			query += "DESC LIMIT $2 OFFSET $3"
		}
		rows, err := conn.Query(ctx, query, param.Author, param.Limit, param.Offset)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		for rows.Next() {
			tags := make([]string, 0)
			var thread models.Thread
			if err := rows.Scan(&thread.Id, &thread.AuthorID, &thread.Title, &thread.Content, &thread.IsLocked, &thread.CreatedAt, &thread.UpdatedAt, &tags); err != nil {
				return nil, 0, err
			}
			thread.Tags = tags
			threads = append(threads, thread)
		}
		return threads, total, nil

	}

	if param.Tag != "" {
		query := `SELECT
			threads.id,
			threads.author_id,
			threads.title,
			threads.content,
			threads.is_locked,
			threads.created_at,
			threads.updated_at,
			array_agg(thread_tags.tag) as tags
		FROM threads
		LEFT JOIN thread_tags
		ON threads.id = thread_tags.thread_id
		WHERE thread_tags.tag = $1
		GROUP BY
			threads.id
		ORDER BY created_at `
		if param.Sort == "old" {
			query += "ASC LIMIT $2 OFFSET $3"
		} else {
			query += "DESC LIMIT $2 OFFSET $3"
		}
		rows, err := conn.Query(ctx, query, param.Tag, param.Limit, param.Offset)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		for rows.Next() {
			tags := make([]string, 0)
			var thread models.Thread
			if err := rows.Scan(&thread.Id, &thread.AuthorID, &thread.Title, &thread.Content, &thread.IsLocked, &thread.CreatedAt, &thread.UpdatedAt, &tags); err != nil {
				return nil, 0, err
			}
			thread.Tags = tags
			threads = append(threads, thread)
		}
		return threads, total, nil
	}

	return nil, 0, fmt.Errorf("Wrong parameters")
}

func GetThreadById(id int) (models.Thread, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return models.Thread{}, fmt.Errorf("Connection error %s", err)
	}
	defer conn.Close()

	tags := make([]string, 0)
	var thread models.Thread
	if err := conn.QueryRow(ctx, `SELECT threads.id,
			threads.author_id,
			threads.title,
			threads.content,
			threads.is_locked,
			threads.created_at,
			threads.updated_at,
			array_agg(thread_tags.tag) as tags
			FROM threads LEFT JOIN thread_tags ON threads.id = thread_tags.thread_id
			WHERE threads.id = $1
			GROUP BY threads.id`,
		id).Scan(&thread.Id, &thread.AuthorID, &thread.Title, &thread.Content, &thread.IsLocked, &thread.CreatedAt, &thread.UpdatedAt, &tags); err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func DeleteThread(userID string, id int) error {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return fmt.Errorf("Connection error %s", err)
	}
	defer conn.Close()

	var threadOwner string
	if err := conn.QueryRow(ctx, `SELECT author_id FROM threads WHERE id = $1`, id).Scan(&threadOwner); err != nil {
		return err
	}
	if threadOwner != userID {
		return fmt.Errorf("This user can not delete this thread")
	}

	if _, err := conn.Exec(ctx, `DELETE FROM threads WHERE id = $1`, id); err != nil {
		return err
	}
	if _, err := conn.Exec(ctx, `DELETE FROM idempotency_keys WHERE user_id = $1 AND response_body -> 'Id' = $2`, userID, id); err != nil {
		return err
	}
	return nil
}

func UpdateThread(newThread models.Thread, userID string, id int) (models.Thread, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return models.Thread{}, fmt.Errorf("Connection error %s", err)
	}
	defer conn.Close()

	var threadOwner string
	var isLocked bool
	var oldTitle string
	var oldContent string
	if err := conn.QueryRow(ctx, `SELECT author_id, is_locked, title, content FROM threads WHERE id = $1`, id).Scan(&threadOwner, &isLocked, &oldTitle, &oldContent); err != nil {
		return models.Thread{}, err
	}

	if threadOwner != userID || isLocked {
		return models.Thread{}, fmt.Errorf("Wrong user or thread is locked")
	}

	if newThread.Content == "" {
		newThread.Content = oldContent
	}
	if newThread.Title == "" {
		newThread.Title = oldTitle
	}

	newThread.Id = id
	if err := conn.QueryRow(ctx, `UPDATE threads SET
		title = $1,
		content = $2,
		updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`, newThread.Title, newThread.Content, id).Scan(&newThread.UpdatedAt); err != nil {
		return models.Thread{}, err
	}

	return newThread, nil
}

func LockThread(lock bool, userID string, id int) error {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return fmt.Errorf("Connection error %s", err)
	}
	defer conn.Close()

	var threadOwner string

	if err := conn.QueryRow(ctx, `SELECT author_id FROM threads WHERE id = $1`, id).Scan(&threadOwner); err != nil {
		return err
	}

	if threadOwner != userID {
		return fmt.Errorf("Wrong user")
	}
	if _, err := conn.Exec(ctx, `UPDATE threads SET
	is_locked = $1 WHERE id = $2`, lock, id); err != nil {
		return err
	}
	return nil
}

func GetSortThreads(sort string) ([]models.Thread, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return nil, fmt.Errorf("Connection error %s", err)
	}
	defer conn.Close()

	threads := make([]models.Thread, 0)
	query := `SELECT threads.id,
			threads.author_id,
			threads.title,
			threads.content,
			threads.is_locked,
			threads.created_at,
			threads.updated_at,
			array_agg(thread_tags.tag) as tags
			FROM threads LEFT JOIN thread_tags ON threads.id = thread_tags.thread_id
			GROUP BY threads.id 
			ORDER BY created_at `
	if sort == "old" {
		query += "ASC"
	} else {
		query += "DESC"
	}

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tags := make([]string, 0)
		var thread models.Thread
		if err := rows.Scan(&thread.Id, &thread.AuthorID, &thread.Title, &thread.Content, &thread.IsLocked, &thread.CreatedAt, &thread.UpdatedAt, &tags); err != nil {
			return nil, err
		}
		thread.Tags = tags
		threads = append(threads, thread)
	}
	return threads, nil
}
