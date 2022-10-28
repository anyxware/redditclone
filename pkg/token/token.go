package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	tokenTTL = 12 * time.Hour
)

type Signer struct {
	signingKey string
}

func NewSigner(signingKey string) Signer {
	return Signer{signingKey: signingKey}
}

type AuthUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type tokenClaims struct {
	User AuthUser `json:"user"`
	jwt.StandardClaims
}

func (s Signer) CreateToken(usr AuthUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		AuthUser{
			Username: usr.Username,
			ID:       usr.ID,
		},

		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		},
	})

	return token.SignedString([]byte(s.signingKey))
}

func (s Signer) ParseToken(accessToken string) (AuthUser, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			return []byte(s.signingKey), nil
		})
	if err != nil {
		return AuthUser{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return AuthUser{}, errors.New("token claims are not of type")
	}

	return claims.User, nil
}
