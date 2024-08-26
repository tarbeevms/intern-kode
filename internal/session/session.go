package session

import (
	"database/sql"

	jwt "github.com/golang-jwt/jwt/v4"
)

type Session struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type JwtCustomClaims struct {
	Name string
	jwt.StandardClaims
}

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepo(conn *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: conn,
	}
}
