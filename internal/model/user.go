package model

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID string `json:"id"`
	Credential
}
