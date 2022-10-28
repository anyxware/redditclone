package slicerepo

import (
	"redditclone/internal/model"
	"redditclone/internal/model/customerr"
	"sync"
)

type usersRepo struct {
	mutex sync.RWMutex
	users []model.User
}

func NewUsersRepo() *usersRepo {
	return &usersRepo{
		users: make([]model.User, 0),
	}
}

func (r *usersRepo) AddUser(user model.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, existedUser := range r.users {
		if existedUser.Username == user.Username {
			return customerr.UserAlreadyExists{Username: user.Username}
		}
	}

	r.users = append(r.users, user)

	return nil
}

func (r *usersRepo) GetUser(cred model.Credential) (model.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, existedUser := range r.users {
		if existedUser.Username == cred.Username && existedUser.Password == cred.Password {
			return existedUser, nil
		}
	}

	return model.User{}, customerr.WrongCredential{Username: cred.Username}
}

func (r *usersRepo) GetUserByID(userID string) (model.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, usr := range r.users {
		if usr.ID == userID {
			return usr, nil
		}
	}

	return model.User{}, customerr.UserNotFoundByID{UserID: userID}
}
