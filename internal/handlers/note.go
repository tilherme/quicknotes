package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"text/template"
)

type noteHandle struct{}

func NewNoteHandle() *noteHandle {
	return &noteHandle{}
}
func (nh *noteHandle) NoteList(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// w.Header().Set("Content-Type", "text/html;charset=utf-8") // setando para json
	// w.Header().Add("Teste", "teste")                   // add um cabeçalho
	// w.Header()["Date"] = nil                           // remover esse cabeçalho"
	// w.Header().Del("Teste"])
	files := []string{
		"../../views/templates/base.html",
		"../../views/templates/pages/home.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Erro: "+err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("execute /")

	t.ExecuteTemplate(w, "base", nil)

}

func (nh *noteHandle) NoteNew(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html;charset=utf-8") // setando para json
	// w.Header().Add("Teste", "teste")                   // add um cabeçalho
	// w.Header()["Date"] = nil                           // remover esse cabeçalho"
	// w.Header().Del("Teste"])
	files := []string{
		"../../views/templates/base.html",
		"../../views/templates/pages/note-new.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Erro: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "base", nil)

}

func (nh *noteHandle) NoteView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Nota não encontrada", http.StatusNotFound)
		return
	}

	files := []string{
		"../../views/templates/base.html",
		"../../views/templates/pages/note-view.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Erro: "+err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "base", id)

}

func (nh *noteHandle) NoteCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(405) //  só pode ser chamado 1 vez por request
		// fmt.Fprint(w, "Metodo não permitido") // opcional o corpo pode ir vazio
		http.Error(w, "Metodo não permitido", http.StatusMethodNotAllowed) // substitui a mensagem e o code
		return                                                             // importante retornar a request caso contrario go vai tentar escrever a resposta
	}
	fmt.Fprint(w, "Criando uma nota")
}
