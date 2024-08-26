package notes

import (
	"time"

	"github.com/google/uuid"
)

// Получение заметки по ID
func (repo *NoteRepository) GetByID(id uuid.UUID) (*Note, error) {
	note := &Note{}
	query := "SELECT id, text, author, created_at, updated_at FROM notes WHERE id = ?"
	err := repo.db.QueryRow(query, id.String()).Scan(&note.ID, &note.Text, &note.Author, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return note, nil
}

// Создание новой заметки
func (repo *NoteRepository) Create(note *Note) error {
	query := "INSERT INTO notes (id, text, author, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	note.ID = uuid.New()
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	_, err := repo.db.Exec(query, note.ID.String(), note.Text, note.Author, note.CreatedAt, note.UpdatedAt)
	return err
}

// Обновление существующей заметки
func (repo *NoteRepository) Update(note *Note) error {
	query := "UPDATE notes SET text = ?, updated_at = ? WHERE id = ?"
	note.UpdatedAt = time.Now()
	_, err := repo.db.Exec(query, note.Text, note.UpdatedAt, note.ID.String())
	return err
}

// Удаление заметки по ID
func (repo *NoteRepository) Delete(id uuid.UUID) error {
	query := "DELETE FROM notes WHERE id = ?"
	_, err := repo.db.Exec(query, id.String())
	return err
}

// Получение всех заметок с данного пользователя
func (repo *NoteRepository) GetAll(username string) ([]*Note, error) {
	notes := []*Note{}

	query := "SELECT id, text, author, created_at, updated_at FROM notes WHERE author = ? ORDER BY updated_at"

	rows, err := repo.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		note := &Note{}
		err := rows.Scan(&note.ID, &note.Text, &note.Author, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
