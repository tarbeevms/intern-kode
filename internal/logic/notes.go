package logic

import (
	"fmt"
	"time"

	"myapp/internal/notes"

	"github.com/google/uuid"
)

// type NoteService interface {
// 	GetNoteByID(id uuid.UUID) (*notes.Note, error)
// 	CreateNote(text string) (*notes.Note, error)
// 	UpdateNote(id uuid.UUID, text string) (*notes.Note, error)
// 	DeleteNote(id uuid.UUID) error
// 	GetNotes(orderBy string) ([]*notes.Note, error)
// }

// type noteService struct {
// 	noteRepo notes.NoteRepository
// }

// func NewNoteService(noteRepo notes.NoteRepository) NoteService {
// 	return &noteService{noteRepo}
// }

func (ll *LogicLayer) CreateNote(text, token string) (*notes.Note, error) {
	correctedText, err := ll.CheckSpelling(text)
	if err != nil {
		return nil, err
	}
	username, err := ll.GetUsernameByToken(token)
	if err != nil {
		return nil, err
	}
	note := &notes.Note{
		ID:        uuid.New(),
		Text:      correctedText,
		Author:    username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := ll.NoteRepo.Create(note); err != nil {
		return nil, err
	}
	return note, nil
}

func (ll *LogicLayer) UpdateNote(id uuid.UUID, text, token string) (*notes.Note, error) {
	note, err := ll.NoteRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	username, err := ll.GetUsernameByToken(token)
	if err != nil {
		return nil, err
	}
	if note.Author != username {
		return nil, fmt.Errorf("no permission to update note")
	}
	correctedText, err := ll.CheckSpelling(text)
	if err != nil {
		return nil, err
	}
	note.Text = correctedText
	note.UpdatedAt = time.Now()
	err = ll.NoteRepo.Update(note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (ll *LogicLayer) DeleteNote(id uuid.UUID, token string) error {
	note, err := ll.NoteRepo.GetByID(id)
	if err != nil {
		return err
	}
	username, err := ll.GetUsernameByToken(token)
	if err != nil {
		return err
	}
	if note.Author != username {
		return fmt.Errorf("no permission to delete note")
	}
	return ll.NoteRepo.Delete(id)
}

func (ll *LogicLayer) GetNotes(token string) ([]*notes.Note, string, error) {
	username, err := ll.GetUsernameByToken(token)
	if err != nil {
		return nil, "", err
	}
	notes, err := ll.NoteRepo.GetAll(username)
	return notes, username, err
}
