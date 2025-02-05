package handlers

import (
	"fmt"

	"gitgub.com/tilherme/quicknotes/internal/models"
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
}

func newRequestNote() (req NoteRequest) {
	req.Color = "color3"
	for i := 1; i <= 9; i++ {
		req.Colors = append(req.Colors, fmt.Sprintf("color%d", i))
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
