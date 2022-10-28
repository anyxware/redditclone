package model

type Author struct {
	ID       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
}
