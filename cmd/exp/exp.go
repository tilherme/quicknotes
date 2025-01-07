package main

import (
	"context"
	"fmt"

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
	// insertPostWithReturn(conn)
	// createTable(conn)
	// insertPost(conn)
	selectById(conn)
}

func createTable(conn *pgx.Conn) {
	query := `
		CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT,
		author TEXT NOT NULL
		);
	`

	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		panic(err)

	}

	fmt.Println("Table 'posts' created")
}

func insertPost(conn *pgx.Conn) {
	title := "titulo post 2"
	content := "conteudo post 2"
	author := "mateus"
	query := `
	INSERT INTO posts(title, content, author)
	values($1 ,$2, $3)
	`
	_, err := conn.Exec(context.Background(), query, title, content, author)
	if err != nil {
		panic(err)
	}
	fmt.Println("Table 'posts' insert")

}

func insertPostWithReturn(conn *pgx.Conn) {
	title := "titulo post guigui"
	content := "conteudo post guigui"
	author := "mateus guigui"
	query := `
	INSERT INTO posts(title, content, author)
	values($1 ,$2, $3) RETURNING id;
	`
	row := conn.QueryRow(context.Background(), query, title, content, author)
	var id int
	err := row.Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("post insert id= ", id)

}

func selectById(conn *pgx.Conn) {
	id := 9
	var title, content, author string
	query := `
	select title, content, author from posts where id = $1;
	`
	row := conn.QueryRow(context.Background(), query, id)
	err := row.Scan(&title, &content, &author)
	if err == pgx.ErrNoRows {
		fmt.Println("Post not found for id", id)
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("post title= %s content= %s author= %s ", title, content, author)

}
