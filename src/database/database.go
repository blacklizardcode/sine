package database

import (
	"github.com/jackc/pgx/v5"
	"os"
	"context"
	"fmt"
)

var DB *pgx.Conn

func InitDB() error {
	var err error
	connStr := "postgres://root:root@127.0.0.1:5432/sine-db"
	DB, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}


	return nil
}