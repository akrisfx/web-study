package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToDB() *pgxpool.Pool {
	// Укажите строку подключения к вашей базе данных
	connString := "postgres://myuser:password@localhost:5432/mydb"

	// config := pgx.ConnConfig{Host: "localhost", Port: 5432, User: "admin", Database: "postgres_db"}
	// config_pool := pgxpool.Config{ConnConfig: }

	// Создайте пул подключений
	pool, err := pgxpool.New(context.Background(), connString)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	// conn, err := pgx.Connect(config)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer conn.Close()
	// defer conn.Close()

	fmt.Println("Connected to the database!")
	return pool
}
