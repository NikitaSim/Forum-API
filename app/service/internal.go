package service

import (
	"context"
	"fmt"

	"stepik.leoscode.http/app/models"

	"github.com/jackc/pgx/v5"
)

func ClearTables() error {

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, models.PostgresqlConnString)
	if err != nil {
		fmt.Println("Connection error", err)
		return err
	}
	defer conn.Close(ctx)

	for _, table := range models.DBTables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)

		_, err := conn.Exec(ctx, query)
		if err != nil {
			fmt.Println("Error truncating table:", table, err)
			return err
		}
	}
	return nil
}
