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
	userRepo := repositories.NewUser(dbpool)
	staticHandler := http.FileServer(http.Dir("views/static"))

	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	noteHandle := handlers.NewNoteHandle(noteRepo)
	userhandle := handlers.NewUserHandler(userRepo)

	mux.Handle("/", handlers.HandleWithError(noteHandle.NoteList))
	mux.Handle("/note/new", handlers.HandleWithError(noteHandle.NoteNew))
	mux.Handle("/note/view", handlers.HandleWithError(noteHandle.NoteView))
	mux.Handle("/note/delete", handlers.HandleWithError(noteHandle.NoteDelete))
	mux.Handle("/note/edit", handlers.HandleWithError(noteHandle.NoteEdit))
	mux.Handle("/note/save", handlers.HandleWithError(noteHandle.NoteSave))

	mux.Handle("GET /user/signup", handlers.HandleWithError(userhandle.SignupForm))
	mux.Handle("POST /user/signup", handlers.HandleWithError(userhandle.Signup))
	mux.Handle("GET /confirmation/{token}", handlers.HandleWithError(userhandle.Confirm))

	mux.Handle("GET /me", handlers.HandleWithError(userhandle.Me))

	mux.Handle("GET /user/signin", handlers.HandleWithError(userhandle.SigninForm))
	mux.Handle("POST /user/signin", handlers.HandleWithError(userhandle.Signin))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), mux); err != nil {
		panic(err)
	}

}
