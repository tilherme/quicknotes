package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"gitgub.com/tilherme/quicknotes/internal/customerror"
	"gitgub.com/tilherme/quicknotes/internal/repositories"
)

type noteHandle struct {
	repo repositories.NoteRepo
}

func NewNoteHandle(repo repositories.NoteRepo) *noteHandle {
	return &noteHandle{repo: repo}
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
		"./views/templates/base.html",
		"./views/templates/pages/home.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Erro: "+err.Error(), http.StatusInternalServerError)
		return
	}

	notes, err := nh.repo.List(r.Context())

	if err != nil {
		fmt.Println(err)
	}

	t.ExecuteTemplate(w, "base", newResponseNoteList(notes))

}

func (nh *noteHandle) NoteNew(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./views/templates/base.html",
		"./views/templates/pages/note-new.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Erro: "+err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "base", newRequestNote())

}

func (nh *noteHandle) NoteView(w http.ResponseWriter, r *http.Request) error {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		return customerror.WithStatus(errors.New("anotação é obrigatoria"), http.StatusBadRequest)
	}
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return err
	}
	files := []string{
		"./views/templates/base.html",
		"./views/templates/pages/note-view.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		return errors.New("aconteceu um erro ao executar essa pagina")
	}
	// ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond)
	// defer cancel()
	note, err := nh.repo.GetById(r.Context(), id)
	if err != nil {
		return err
	}
	buff := &bytes.Buffer{}
	err = t.ExecuteTemplate(buff, "base", newResponseNote(note))
	if err != nil {
		return err
	}
	buff.WriteTo(w)
	return nil

}

func (nh *noteHandle) NoteCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(405) //  só pode ser chamado 1 vez por request
		// fmt.Fprint(w, "Metodo não permitido") // opcional o corpo pode ir vazio
		// errors.New("deu ruim")
		http.Error(w, "Metodo não permitido", http.StatusMethodNotAllowed) // substitui a mensagem e o code
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	color := r.PostForm.Get("color")

	note, err := nh.repo.Create(r.Context(), title, content, color)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, fmt.Sprintf("/note/view?id=%d", note.Id.Int), http.StatusSeeOther)
}

func (nh *noteHandle) NoteDelete(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodDelete)
		// w.WriteHeader(405) //  só pode ser chamado 1 vez por request
		// fmt.Fprint(w, "Metodo não permitido") // opcional o corpo pode ir vazio
		// errors.New("deu ruim")
		http.Error(w, "Metodo não permitido", http.StatusMethodNotAllowed) // substitui a mensagem e o code
	}
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		return customerror.WithStatus(errors.New("anotação é obrigatoria"), http.StatusBadRequest)
	}
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return err
	}
	err = nh.repo.Delete(r.Context(), id)
	if err != nil {
		return err
	}
	return nil
}
