package service

import (
	"context"
	"fmt"

	"stepik.leoscode.http/app/models"

	"github.com/jackc/pgx/v5"
)

func LoginUser(user models.Users) (string, error) {

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, models.PostgresqlConnString)
	if err != nil {
		fmt.Println("Connection error", err)
		return "", err
	}
	defer conn.Close(ctx)

	var id string
	query := `INSERT INTO users (id, username, password)
		VALUES (gen_random_uuid(), $1, $2)
		ON CONFLICT (username) DO NOTHING
		RETURNING id;`
	row := conn.QueryRow(ctx, query, user.Name, user.Password)

	if err = row.Scan(&id); err != nil {
		row = conn.QueryRow(ctx, "SELECT id FROM users WHERE username = $1", user.Name)
		err = row.Scan(&id)
		if err != nil {
			return "", err
		}
	}

	return id, nil
}
