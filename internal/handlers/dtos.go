package handlers

import (
	"fmt"

	"gitgub.com/tilherme/quicknotes/internal/models"
	"gitgub.com/tilherme/quicknotes/internal/validators"
)

type NoteResponse struct {
	Id      int
	Title   string
	Content string
	Color   string
}
type NoteRequest struct {
	Id      int
	Title   string
	Content string
	Color   string
	Colors  []string
	validators.FormValidation
}
type UserRequest struct {
	Email    string
	Password string
	Name     string
	validators.FormValidation
}

func newUserRequest(email, password, name string) (req UserRequest) {
	req.Email = email
	req.Password = password
	req.Name = name
	return
}

func newRequestNote(note *models.Note) (req NoteRequest) {
	for i := 1; i <= 9; i++ {
		req.Colors = append(req.Colors, fmt.Sprintf("color%d", i))
	}
	if note != nil {
		req.Id = int(note.Id.Int.Int64())
		req.Title = note.Title.String
		req.Content = note.Content.String
		req.Color = note.Color.String

	} else {
		req.Color = "color3"
	}
	return
}

func newResponseNote(note *models.Note) (res NoteResponse) {
	res.Id = int(note.Id.Int.Int64())
	res.Title = note.Title.String
	res.Content = note.Content.String
	res.Color = note.Color.String
	return
}

func newResponseNoteList(notes []models.Note) (res []NoteResponse) {

	for _, note := range notes {
		res = append(res, newResponseNote(&note))

	}
	return
}
