package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT,8000"`
	Password   string `env:"PASSWORD,required"`
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
		value := os.Getenv(envName)
		f := reflect.ValueOf(&conf).Elem().FieldByName(field.Name)
		f.SetString(value)
	}
	return
}
func (c Config) Vslidade() {
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
	config.Vslidade()
	// config.ServerPort = os.Getenv("SERVER_PORT")

	// fmt.Println(config.SPrint()
	return config
}
