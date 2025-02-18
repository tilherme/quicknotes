package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("ola")
	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hash))
	err = bcrypt.CompareHashAndPassword(hash, []byte("1234456"))
	if err != nil {
		fmt.Println("Senha invalida")
	} else {
		fmt.Println("Login feito com sucesso!")

	}
}
