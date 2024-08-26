package logic

import (
	"errors"
	"myapp/internal/users"
)

func (ll *LogicLayer) VerifyUserCredentials(username, password string) (*users.LoginRequest, error) {
	// Проверяем учетные данные пользователя
	if username == "" || password == "" {
		return nil, errors.New("invalid username or password")
	}
	user, err := ll.UserRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user.Password != password || username == "" || password == "" {
		return nil, errors.New("invalid username or password")
	}
	return user, nil
}
