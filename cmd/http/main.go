package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"gitgub.com/tilherme/quicknotes/internal/handlers"
)

func main() {
	config := LoadConfig()

	slog.SetDefault(newLogger(os.Stderr, config.GetLevelLog()))
	mux := http.NewServeMux()
	slog.Info(fmt.Sprintf("Senha: %s", config.Password))
	slog.Info(fmt.Sprintf("deu bom porta %s\n", config.ServerPort))
	slog.Info(fmt.Sprintf("teste %s\n", config.Teste))

	staticHandler := http.FileServer(http.Dir("../../views/static"))

	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	noteHandle := handlers.NewNoteHandle()
	mux.HandleFunc("/", noteHandle.NoteList)
	mux.HandleFunc("/note/new", noteHandle.NoteNew)
	mux.HandleFunc("/note/view", noteHandle.NoteNew)
	mux.HandleFunc("/note/create", noteHandle.NoteCreate)

	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), mux)

}
