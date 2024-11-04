package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
}

func LoadConfig() Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}
	config := Config{}
	config.ServerPort = os.Getenv("SERVER_PORT")
	return config
}
