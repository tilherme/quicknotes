package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	fmt.Println(GenerateToken())

}

func GenerateToken() string {
	r := make([]byte, 32)
	rand.Read(r)
	return base64.URLEncoding.EncodeToString(r)

}
