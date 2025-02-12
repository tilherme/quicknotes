package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gitgub.com/tilherme/quicknotes/internal/customerror"
	"gitgub.com/tilherme/quicknotes/internal/models"
	"gitgub.com/tilherme/quicknotes/internal/repositories"
)

type noteHandle struct {
	repo repositories.NoteRepo
}

func NewNoteHandle(repo repositories.NoteRepo) *noteHandle {
	return &noteHandle{repo: repo}
}

func (nh *noteHandle) NoteList(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return errors.ErrUnsupported
	}

	notes, err := nh.repo.List(r.Context())

	if err != nil {
		return err
	}

	return render(w, http.StatusOK, "home.html", newResponseNoteList(notes))

}

func (nh *noteHandle) NoteNew(w http.ResponseWriter, r *http.Request) error {
	return render(w, http.StatusOK, "note-new.html", newRequestNote(nil))

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
	note, err := nh.repo.GetById(r.Context(), id)
	if err != nil {
		return err
	}
	return render(w, http.StatusOK, "note-view.html", newResponseNote(note))

}

func (nh *noteHandle) NoteSave(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Metodo não permitido", http.StatusMethodNotAllowed) // substitui a mensagem e o code
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	idParam := r.PostForm.Get("id")
	id, _ := strconv.Atoi(idParam)
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	color := r.PostForm.Get("color")

	data := newRequestNote(nil)
	data.Id = id
	data.Title = title
	data.Content = content
	data.Color = color
	if strings.TrimSpace(content) == "" {
		data.AddFieldErrors("content", "Campo é obrigatorio")
	}

	if !data.Valid() {
		if id > 0 {
			render(w, http.StatusUnprocessableEntity, "note-edit.html", data)
		} else {
			render(w, http.StatusUnprocessableEntity, "note-new.html", data)
		}
		return nil
	}

	var note *models.Note
	if id > 0 {
		note, err = nh.repo.Update(r.Context(), id, title, content, color)
	} else {
		note, err = nh.repo.Create(r.Context(), title, content, color)

	}

	if err != nil {
		return err
	}
	http.Redirect(w, r, fmt.Sprintf("/note/view?id=%d", note.Id.Int), http.StatusSeeOther)
	return nil
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

func (nh *noteHandle) NoteEdit(w http.ResponseWriter, r *http.Request) error {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		return customerror.WithStatus(errors.New("anotação é obrigatoria"), http.StatusBadRequest)
	}
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return err
	}
	note, err := nh.repo.GetById(r.Context(), id)
	if err != nil {
		return err
	}
	return render(w, http.StatusOK, "note-edit.html", newRequestNote(note))

}
