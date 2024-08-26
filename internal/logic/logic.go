package logic

import (
	"myapp/internal/notes"
	"myapp/internal/session"
	"myapp/internal/users"
)

// LogicLayer предоставляет бизнес-логику для работы с данными.
type LogicLayer struct {
	NoteRepo    *notes.NoteRepository
	UserRepo    *users.UserRepository
	SessionRepo *session.SessionRepository
}

// NewLogicLayer создает новый экземпляр слоя логики.
func NewLogicLayer(nr *notes.NoteRepository, ur *users.UserRepository, sr *session.SessionRepository) *LogicLayer {
	return &LogicLayer{
		NoteRepo:    nr,
		UserRepo:    ur,
		SessionRepo: sr,
	}
}
