package session

import (
	"net/http"
)

type Token struct {
	Value string
}

type AuthUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Manager interface {
	GetTokenFromRequest(r *http.Request) (Token, error)
	GetUserByToken(token Token) (AuthUser, error)
	CreateToken(user AuthUser) (Token, error)
	WriteToken(w http.ResponseWriter, token Token) error
}
