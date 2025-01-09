package main

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT,9000"`
	DbConnUrl  string `env:"DB_CONN_URL,required"`
	LevelLog   string `env:"LEVEL_LOG,info"`
}

func (c Config) GetLevelLog() slog.Level {
	switch strings.ToLower(c.LevelLog) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (c Config) SPrint() (envs string) {
	v := reflect.ValueOf(c)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envTag := strings.Split(field.Tag.Get("env"), ",")

		name := envTag[0]
		value := envTag[1]
		envs += fmt.Sprintf("%s - %s\n", name, value)
	}
	return
}

func (c Config) LoadFromEnv() (conf Config) {
	v := reflect.ValueOf(c)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envTag := strings.Split(field.Tag.Get("env"), ",")

		envName := envTag[0]
		defaultValue := envTag[1]
		value := os.Getenv(envName)
		if value == "" && value != "required" {
			f := reflect.ValueOf(&conf).Elem().FieldByName(field.Name)
			f.SetString(defaultValue)
		} else {
			f := reflect.ValueOf(&conf).Elem().FieldByName(field.Name)
			f.SetString(value)
		}

	}
	return
}
func (c Config) Validade() {
	var validadeMsg string
	v := reflect.ValueOf(c)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		envTag := strings.Split(t.Field(i).Tag.Get("env"), ",")

		name := envTag[0]
		envValue := envTag[1]
		if envValue == "required" && value.String() == "" {
			validadeMsg += fmt.Sprintf("%s is required", name)

		}

	}
	if len(validadeMsg) != 0 {
		panic(validadeMsg)
	}
}

func LoadConfig() Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		// panic(err)
		fmt.Println(err)
	}
	config := Config{}
	config = config.LoadFromEnv()
	config.Validade()
	// config.ServerPort = os.Getenv("SERVER_PORT")

	// fmt.Println(config.SPrint()
	return config
}
