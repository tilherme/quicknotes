package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println(err)
	}
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	fmt.Printf("%s:%s", host, port)
}
