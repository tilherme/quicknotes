package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func main() {
	dbUrl := "postgres://postgres:secret@localhost:5432/postgres"
	con, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Conecxão feita com sucesso!!! ")
	defer con.Close(context.Background())
}
