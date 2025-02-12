package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"text/template"
)

func render(w http.ResponseWriter, status int, page string, data any) error {
	files := []string{
		"views/templates/base.html",
	}
	files = append(files, "views/templates/pages/"+page)
	t, err := template.ParseFiles(files...)
	if err != nil {
		return errors.New("aconteceu um erro ao executar essa pagina")
	}
	if err != nil {
		return err
	}
	buff := &bytes.Buffer{}
	err = t.ExecuteTemplate(buff, "base", data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	buff.WriteTo(w)
	return nil
}
