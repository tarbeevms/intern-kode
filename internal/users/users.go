package users

import (
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepo(conn *sql.DB) *UserRepository {
	return &UserRepository{
		db: conn,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
