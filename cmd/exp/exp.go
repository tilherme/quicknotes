package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
	var err error
	var conn *pgx.Conn

	dbUrl := "postgres://postgres:secret@localhost:5432/postgres"
	conn, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		panic(err)
	}

	defer conn.Close(context.Background())
	fmt.Println("Conecxão feita com sucesso!!! ", conn)

	createTable(conn)
}

func createTable(conn *pgx.Conn) {
	query := `
		CREATE TABLE posts (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT,
		author TEXT NOT NULL
		);
	`

	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Error creating table: %v\n", err)
		return
	}

	fmt.Println("Table 'posts' created")
}
