package session

import (
	"database/sql"
	"fmt"
)

// AddSession добавляет новую сессию в MySQL.
func (sr *SessionRepository) AddSession(sess *Session) error {
	query := "REPLACE INTO sessions (username, token) VALUES (?, ?)"
	_, err := sr.db.Exec(query, sess.Username, sess.Token)
	if err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}
	return nil
}

// DeleteSessionByToken удаляет сессию по токену.
func (sr *SessionRepository) DeleteSessionByToken(token string) error {
	query := "DELETE FROM sessions WHERE token = ?"
	_, err := sr.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// GetSessionByToken получает сессию по токену.
func (sr *SessionRepository) GetSessionByToken(token string) (*Session, error) {
	query := "SELECT username, token FROM sessions WHERE token = ? LIMIT 1"
	row := sr.db.QueryRow(query, token)

	var sess Session
	err := row.Scan(&sess.Username, &sess.Token)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	return &sess, nil
}
