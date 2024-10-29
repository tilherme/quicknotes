package main

import (
	"fmt"
	"html/template"
	"os"
)

type templateData struct {
	Nome string
}

func main() {
	t, err := template.ParseFiles("layout.html", "layout2.html", "header.html", "footer.html")
	if err != nil {
		panic(err)
	}
	fmt.Println(t.Name())
	// fmt.Println(t.DefinedTemplates())

	// data := templateData{Nome: "gui"}
	err = t.ExecuteTemplate(os.Stdout, "layout2.html", "2024")
	if err != nil {
		panic(err)
	}
}
