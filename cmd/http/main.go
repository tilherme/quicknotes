package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"gitgub.com/tilherme/quicknotes/internal/handlers"
	"gitgub.com/tilherme/quicknotes/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config := LoadConfig()

	slog.SetDefault(newLogger(os.Stderr, config.GetLevelLog()))
	dbpool, err := pgxpool.New(context.Background(), config.DbConnUrl)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("Conexão feita com sucesso")
	defer dbpool.Close()

	mux := http.NewServeMux()

	noteRepo := repositories.NewNote(dbpool)
	staticHandler := http.FileServer(http.Dir("views/static"))

	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	noteHandle := handlers.NewNoteHandle(noteRepo)

	mux.HandleFunc("/", noteHandle.NoteList)
	mux.HandleFunc("/note/new", noteHandle.NoteNew)
	mux.Handle("/note/view", handlers.HandleWithError(noteHandle.NoteView))
	mux.Handle("/note/delete", handlers.HandleWithError(noteHandle.NoteDelete))

	mux.HandleFunc("/note/create", noteHandle.NoteCreate)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), mux); err != nil {
		panic(err)
	}

}
