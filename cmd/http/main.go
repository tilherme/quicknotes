package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func noteList(w http.ResponseWriter, r *http.Request) {
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

	t.ExecuteTemplate(w, "base", nil)

}

func noteNew(w http.ResponseWriter, r *http.Request) {
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

func noteView(w http.ResponseWriter, r *http.Request) {
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

func noteCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(405) //  só pode ser chamado 1 vez por request
		// fmt.Fprint(w, "Metodo não permitido") // opcional o corpo pode ir vazio
		http.Error(w, "Metodo não permitido", http.StatusMethodNotAllowed) // substitui a mensagem e o code
		return                                                             // importante retornar a request caso contrario go vai tentar escrever a resposta
	}
	fmt.Fprint(w, "Criando uma nota")
}

func main() {
	config := LoadConfig()
	mux := http.NewServeMux()
	fmt.Printf("%s", config.ServerPort)
	staticHandler := http.FileServer(http.Dir("../../views/static"))

	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	mux.HandleFunc("/", noteList)
	mux.HandleFunc("/note/new", noteNew)
	mux.HandleFunc("/note/view", noteView)
	mux.HandleFunc("/note/create", noteCreate)

	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), mux)

}
