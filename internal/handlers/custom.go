package handlers

import (
	"errors"
	"net/http"
	"text/template"

	"gitgub.com/tilherme/quicknotes/internal/customerror"
)

type HandleWithError func(w http.ResponseWriter, r *http.Request) error

type responseWriterWrapper struct {
	http.ResponseWriter
	wroteHeader bool
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	if !rw.wroteHeader {
		rw.ResponseWriter.WriteHeader(statusCode)
		rw.wroteHeader = true
	}
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (f HandleWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := &responseWriterWrapper{ResponseWriter: w}

	if err := f(rw, r); err != nil {
		var statusErr customerror.StatusError

		if errors.As(err, &statusErr) {
			if statusErr.StatusCode() == http.StatusNotFound {
				files := []string{
					"../../views/templates/base.html",
					"../../views/templates/pages/404.html",
				}

				t, tmplErr := template.ParseFiles(files...)
				if tmplErr != nil {
					if !rw.wroteHeader {
						http.Error(rw, tmplErr.Error(), http.StatusInternalServerError)
					}
					return
				}

				if !rw.wroteHeader {
					rw.WriteHeader(http.StatusNotFound)
				}
				t.ExecuteTemplate(rw, "base", statusErr.Error())
				return
			}

			if !rw.wroteHeader {
				http.Error(rw, err.Error(), statusErr.StatusCode())
			}
			return
		}

		if !rw.wroteHeader {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}
