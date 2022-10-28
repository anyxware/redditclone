package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"redditclone/internal/model"
	"redditclone/internal/model/customerr"
	"redditclone/pkg/hexid"
)

type usersRepo interface {
	AddUser(user model.User) error
	GetUser(cred model.Credential) (model.User, error)
	GetUserByID(userID string) (model.User, error)
}

func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (s *service) RegisterUser(cred model.Credential) (model.User, error) {
	id, err := hexid.Generate()
	if err != nil {
		return model.User{}, err
	}

	usr := model.User{ID: id, Credential: cred}
	usr.Password = hashPassword(usr.Password)

	err = s.usersRepo.AddUser(usr)
	if err == nil {
		logrus.Infof("user registered: %s", cred.Username)
	}

	return usr, err
}

func (s *service) LoginUser(cred model.Credential) (model.User, error) {
	cred.Password = hashPassword(cred.Password)

	usr, err := s.usersRepo.GetUser(cred)
	if _, ok := err.(customerr.WrongCredential); ok {
		return model.User{}, customerr.WrongCredential{Username: cred.Username}
	}
	if err == nil {
		logrus.Infof("user logged: %s", cred.Username)
	}

	return usr, err
}

func (s *service) GetUserByID(userID string) (model.User, error) {
	return s.usersRepo.GetUserByID(userID)
}
