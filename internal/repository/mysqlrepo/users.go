package mysqlrepo

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"redditclone/internal/model"
	"redditclone/internal/model/customerr"
	"reflect"
)

type usersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *usersRepo {
	return &usersRepo{db: db}
}

func (r *usersRepo) AddUser(user model.User) error {
	_, err := r.db.Exec(
		"INSERT INTO user (`id`, `username`, `password`) VALUES (?, ?, ?)",
		user.ID,
		user.Username,
		user.Password,
	)
	// Error 1062: Duplicate entry 'van' for key 'user.username'
	if mysqlErr, ok := err.(*mysql.MySQLError); ok && reflect.DeepEqual(mysqlErr, &mysql.MySQLError{
		Number:  1062,
		Message: fmt.Sprintf("Duplicate entry '%s' for key 'user.username'", user.Username),
	}) {
		return customerr.UserAlreadyExists{Username: user.Username}
	}
	return err
}

func (r *usersRepo) GetUser(cred model.Credential) (model.User, error) {
	var user model.User
	err := r.db.QueryRow(
		"SELECT id, username, password FROM user WHERE username = ? AND password = ?",
		cred.Username,
		cred.Password,
	).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return model.User{}, customerr.WrongCredential{Username: cred.Username}
	}
	return user, err
}

func (r *usersRepo) GetUserByID(userID string) (model.User, error) {
	var user model.User
	err := r.db.QueryRow(
		"SELECT id, username, password FROM user WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return model.User{}, customerr.UserNotFoundByID{UserID: userID}
	}
	return user, err
}
