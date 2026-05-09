package service

import (
	"context"
	"fmt"

	"stepik.leoscode.http/app/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UploudFile(attc models.Attachments, userID string) (models.Attachments, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return models.Attachments{}, fmt.Errorf("Connection error %s", err)
	}
	defer conn.Close()

	if err := conn.QueryRow(ctx, `INSERT INTO attachments
	(id, thread_id, filename, mime_type, size, path) VALUES
	(gen_random_uuid(), $1, $2, $3, $4, $5)
	RETURNING id, created_at`, attc.ThreadId, attc.Filename, attc.MimeType, attc.Size, attc.Path).Scan(&attc.Id, &attc.CreatedAt); err != nil {
		return models.Attachments{}, err
	}

	return attc, nil
}

func ShowFiles(id int) ([]models.Attachments, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return nil, fmt.Errorf("Connection error %s", err)
	}
	defer conn.Close()

	rows, err := conn.Query(ctx, `SELECT id, thread_id, filename, mime_type, size,created_at FROM attachments WHERE thread_id = $1`, id)
	if err != nil {
		return nil, err
	}

	attcs := make([]models.Attachments, 0)
	for rows.Next() {
		var attc models.Attachments
		if err := rows.Scan(&attc.Id, &attc.ThreadId, &attc.Filename, &attc.MimeType, &attc.Size, &attc.CreatedAt); err != nil {
			return nil, err
		}
		attcs = append(attcs, attc)
	}

	return attcs, nil
}

func GetFile(fileID string) (string, string, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, models.PostgresqlConnString)
	if err != nil {
		return "", "", fmt.Errorf("Connection error %s", err)
	}
	defer conn.Close()

	var path, name string
	if err := conn.QueryRow(ctx, `SELECT path, filename FROM attachments WHERE id = $1`, fileID).Scan(&path, &name); err != nil {
		return "", "", err
	}
	return path, name, nil
}
