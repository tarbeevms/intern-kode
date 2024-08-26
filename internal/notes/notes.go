package notes

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteText struct {
	Text string `json:"text"`
}

type NoteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db}
}
