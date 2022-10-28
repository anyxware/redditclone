package cookie

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"redditclone/pkg/session"
)

const (
	cookieName = "RedditcloneCookie"
)

type storage interface {
	Add(mkey string, serialized []byte) error
	Get(mkey string) ([]byte, error)
	GetTTL() int
}

type manager struct {
	storage storage
	ttl     int
}

func NewManager(storage storage) *manager {
	return &manager{storage: storage, ttl: storage.GetTTL()}
}

func (m *manager) GetTokenFromRequest(r *http.Request) (session.Token, error) {
	var cookie *http.Cookie
	for _, c := range r.Cookies() {
		cookie = c
		break
	}
	if cookie == nil {
		return session.Token{}, errors.New(fmt.Sprintf("cookie with name '%s' not found", cookieName))
	}
	return session.Token{Value: cookie.Value}, nil
}

func (m *manager) GetUserByToken(token session.Token) (session.AuthUser, error) {
	serialized, err := m.storage.Get(token.Value)
	if err != nil {
		return session.AuthUser{}, err
	}

	var user session.AuthUser
	err = json.Unmarshal(serialized, &user)

	return user, err
}

func (m *manager) CreateToken(user session.AuthUser) (session.Token, error) {
	idHash := md5.Sum([]byte(user.ID))
	usernameHash := md5.Sum([]byte(user.ID))
	wholeHash := md5.Sum(append(idHash[:], usernameHash[:]...))

	mkey := base64.StdEncoding.EncodeToString(wholeHash[:])
	serialized, err := json.Marshal(user)
	if err != nil {
		return session.Token{}, err
	}

	err = m.storage.Add(mkey, serialized)
	if err != nil {
		return session.Token{}, err
	}

	return session.Token{Value: mkey}, nil
}

func (m *manager) WriteToken(w http.ResponseWriter, token session.Token) error {
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  token.Value,
		MaxAge: m.ttl,
	}
	http.SetCookie(w, cookie)
	return nil
}
