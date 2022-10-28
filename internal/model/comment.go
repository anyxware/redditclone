package model

import "time"

type Comment struct {
	ID      string `json:"id" bson:"id"`
	Created string `json:"created" bson:"created"`
	Author  Author `json:"author" bson:"author"`
	Body    string `json:"body" bson:"body"`
}

func NewComment(commentID string, text string, author Author) Comment {
	return Comment{
		ID:      commentID,
		Created: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		Author:  author,
		Body:    text,
	}
}
