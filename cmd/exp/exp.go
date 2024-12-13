package main

import (
	"errors"
	"fmt"
	"os"
)

type CustomError struct {
	msg  string
	code int
}

var bad = &CustomError{msg: "vai dar merda", code: 400}

func (c *CustomError) Error() string {
	return fmt.Sprintf("%s, %d", c.msg, c.code)

}
func NewCustomeError(msg string, code int) error {
	return &CustomError{msg: msg, code: code}
}
func process() (string, error) {
	f, err := os.Open("../../config.son")
	if err != nil {
		// return "", errors.New("erro new")
		// return "", fmt.Errorf("vai dar não")
		// return "", NewCustomeError("vai dar merda", 400)
		// return "", bad
		return "", fmt.Errorf("vai dar não; (%w)", err)

	}
	return f.Name(), nil
}
func main() {
	r, err := process()
	if err != nil {
		// var bad *CustomError
		// if errors.As(err, &bad) {
		// 	fmt.Println(bad.code)
		// 	return
		// }
		// fmt.Println(err)
		fmt.Println(errors.Unwrap(err))
		return
	}
	fmt.Println(r, "fff")
}
