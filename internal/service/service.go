package service

import "sync"

type service struct {
	usersRepo  usersRepo
	postsMutex sync.Mutex
	postsRepo  postsRepo
}

func NewService(usersRepo usersRepo, postsRepo postsRepo) *service {
	return &service{
		usersRepo: usersRepo,
		postsRepo: postsRepo,
	}
}
