package handlers

import (
	"errors"
	"net/http"
	"text/template"

	"gitgub.com/tilherme/quicknotes/internal/customerror"
)

type HandleWithError func(w http.ResponseWriter, r *http.Request) error

func (f HandleWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		var statusErr customerror.StatusError
		if errors.As(err, &statusErr) {
			if statusErr.StatusCode() == http.StatusNotFound {
				files := []string{
					"../../views/templates/base.html",
					"../../views/templates/pages/404.html",
				}
				t, err := template.ParseFiles(files...)
				if err != nil {
					http.Error(w, err.Error(), statusErr.StatusCode())
				}
				t.ExecuteTemplate(w, "base", statusErr.Error())
				return
			}
			http.Error(w, err.Error(), statusErr.StatusCode())
			return

		}
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
}
