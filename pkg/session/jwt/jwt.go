package jwtsession

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"redditclone/pkg/session"
	"strings"
	"time"
)

const (
	authorizationHeader = "Authorization"
)

type tokenClaims struct {
	User session.AuthUser `json:"user"`
	jwt.StandardClaims
}

type manager struct {
	signingKey string
	ttl        int
}

func NewManager(signingKey string, ttl int) *manager {
	return &manager{signingKey: signingKey, ttl: ttl}
}

func (m *manager) GetTokenFromRequest(r *http.Request) (session.Token, error) {
	header := r.Header.Get(authorizationHeader)
	if header == "" {
		return session.Token{}, errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return session.Token{}, errors.New("invalid auth header")
	}

	return session.Token{Value: headerParts[1]}, nil
}

func (m *manager) GetUserByToken(token session.Token) (session.AuthUser, error) {
	accessToken, err := jwt.ParseWithClaims(
		token.Value,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			return []byte(m.signingKey), nil
		})
	if err != nil {
		return session.AuthUser{}, err
	}

	claims, ok := accessToken.Claims.(*tokenClaims)
	if !ok {
		return session.AuthUser{}, errors.New("token claims are not of type")
	}

	return claims.User, nil
}

func (m *manager) CreateToken(user session.AuthUser) (session.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		session.AuthUser{
			Username: user.Username,
			ID:       user.ID,
		},

		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Duration(m.ttl) * time.Second).Unix(),
		},
	})

	accessToken, err := token.SignedString([]byte(m.signingKey))

	return session.Token{Value: accessToken}, err
}

func (m *manager) WriteToken(w http.ResponseWriter, token session.Token) error {
	resp := []byte(
		fmt.Sprintf("{\"token\": \"%s\"}", token.Value),
	)

	if _, err := w.Write(resp); err != nil {
		return err
	}

	return nil
}
