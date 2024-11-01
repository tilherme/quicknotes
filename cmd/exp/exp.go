package main

import (
	"flag"
	"fmt"
)

func main() {
	var port string
	var verbose bool
	var value int64
	flag.StringVar(&port, "port", "7000", "Server port")
	flag.BoolVar(&verbose, "v", false, "Verbose")
	flag.Int64Var(&value, "value", 0, "value sum")
	flag.Parse()
	if verbose {
		fmt.Println("SERVER: ", port)
		fmt.Println("Valor ", value)

	} else {
		fmt.Println("SERVER ERRO: ")

	}
}
