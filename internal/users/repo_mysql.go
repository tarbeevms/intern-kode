package users

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetUserByUsername ищет пользователя по username в MySQL
func (ur *UserRepository) GetUserByUsername(username string) (*LoginRequest, error) {
	query := "SELECT username, password FROM users WHERE username = ? LIMIT 1"
	row := ur.db.QueryRow(query, username)

	var user LoginRequest
	err := row.Scan(&user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}
