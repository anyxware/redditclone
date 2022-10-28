package customerr

import (
	"fmt"
)

type UserAlreadyExists struct {
	Username string
}

func (e UserAlreadyExists) Error() string {
	return fmt.Sprintf("user already exists: %s", e.Username)
}

type WrongCredential struct {
	Username string
}

func (e WrongCredential) Error() string {
	return fmt.Sprintf("wrong credential for user: %s", e.Username)
}

type UserNotFoundByID struct {
	UserID string
}

func (e UserNotFoundByID) Error() string {
	return fmt.Sprintf("user not found by ID: %s", e.UserID)
}

type UserNotFoundByUsername struct {
	Username string
}

func (e UserNotFoundByUsername) Error() string {
	return fmt.Sprintf("user not found by username: %s", e.Username)
}

type PostNotFoundByID struct {
	PostID string
}

func (e PostNotFoundByID) Error() string {
	return fmt.Sprintf("post not found by ID: %s", e.PostID)
}

type CommentNotFoundByID struct {
	PostID    string
	CommentID string
}

func (e CommentNotFoundByID) Error() string {
	return fmt.Sprintf("comment not found by ID: %s in post with ID: %s", e.CommentID, e.PostID)
}

type NotOwner struct {
	Username string
}

func (e NotOwner) Error() string {
	return fmt.Sprintf("user %s not own this resource", e.Username)
}

type RequestNotParsed struct {
	Message string
}

func (e RequestNotParsed) Error() string {
	return fmt.Sprintf("request wasn't parsed: %s", e.Message)
}

type Unauthorized struct {
	Message string
}

func (e Unauthorized) Error() string {
	return fmt.Sprintf("user unauthorized: %s", e.Message)
}
