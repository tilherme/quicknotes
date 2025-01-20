package handlers

import (
	"gitgub.com/tilherme/quicknotes/internal/models"
)

type NoteResponse struct {
	Id      int
	Title   string
	Content string
}

func newResponseNote(note *models.Note) (res NoteResponse) {
	res.Id = int(note.Id.Int.Int64())
	res.Title = note.Title.String
	res.Content = note.Content.String
	return
}

func newResponseNoteList(notes []models.Note) (res []NoteResponse) {

	for _, note := range notes {
		res = append(res, newResponseNote(&note))

	}
	return
}
