package mysqlrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"redditclone/internal/model"
	"redditclone/internal/model/customerr"
	"reflect"
	"testing"
)

func compareErrorsMsg(err1 error, err2 error) bool {
	return fmt.Sprint(err1) == fmt.Sprint(err2)
}

func TestAddUser(t *testing.T) {
	cases := []struct {
		user        model.User
		expectedErr error
		run         func(user model.User, repo *usersRepo, mock sqlmock.Sqlmock) error
	}{
		{
			user:        model.User{ID: "1"},
			expectedErr: nil,
			run: func(user model.User, repo *usersRepo, mock sqlmock.Sqlmock) error {
				mock.
					ExpectExec("INSERT INTO user").
					WithArgs(user.ID, user.Username, user.Password).
					WillReturnResult(sqlmock.NewResult(1, 1))
				return repo.AddUser(user)
			},
		},
		{
			user:        model.User{ID: "2"},
			expectedErr: errors.New("bad query"),
			run: func(user model.User, repo *usersRepo, mock sqlmock.Sqlmock) error {
				mock.
					ExpectExec("INSERT INTO user").
					WithArgs(user.ID, user.Username, user.Password).
					WillReturnError(errors.New("bad query"))
				return repo.AddUser(user)
			},
		},
		{
			user:        model.User{ID: "3", Credential: model.Credential{Username: "ivan"}},
			expectedErr: customerr.UserAlreadyExists{Username: "ivan"},
			run: func(user model.User, repo *usersRepo, mock sqlmock.Sqlmock) error {
				mock.
					ExpectExec("INSERT INTO user").
					WithArgs(user.ID, user.Username, user.Password).
					WillReturnError(&mysql.MySQLError{
						Number:  1062,
						Message: fmt.Sprintf("Duplicate entry 'ivan' for key 'user.username'"),
					})
				return repo.AddUser(user)
			},
		},
	}

	for i, item := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("cant create mock: %s", err)
		}

		repo := NewUsersRepo(db)
		err = item.run(item.user, repo, mock)
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
		db.Close()
	}
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUsersRepo(db)

	cases := []struct {
		expectedUser model.User
		expectedErr  error
		run          func(user model.User) (model.User, error)
	}{
		{
			expectedUser: model.User{ID: "1", Credential: model.Credential{Username: "ivan", Password: "qqq"}},
			expectedErr:  nil,
			run: func(user model.User) (model.User, error) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"})
				rows.AddRow(user.ID, user.Username, user.Password)
				mock.
					ExpectQuery("SELECT id, username, password FROM user WHERE").
					WithArgs(user.Username, user.Password).
					WillReturnRows(rows)
				return repo.GetUser(user.Credential)
			},
		},
		{
			expectedUser: model.User{},
			expectedErr:  errors.New("bad query"),
			run: func(user model.User) (model.User, error) {
				mock.
					ExpectQuery("	SELECT id, username, password FROM user WHERE").
					WithArgs(user.Username, user.Password).
					WillReturnError(errors.New("bad query"))
				return repo.GetUser(user.Credential)
			},
		},
		{
			expectedUser: model.User{},
			expectedErr:  customerr.WrongCredential{Username: "ivan"},
			run: func(user model.User) (model.User, error) {
				user.Username = "ivan"
				mock.
					ExpectQuery("SELECT id, username, password FROM user WHERE").
					WithArgs(user.Username, user.Password).
					WillReturnError(sql.ErrNoRows)
				return repo.GetUser(user.Credential)
			},
		},
	}

	for i, item := range cases {
		user, err := item.run(item.expectedUser)
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
		if !reflect.DeepEqual(item.expectedUser, user) {
			t.Errorf("[%d] expected user: %+v, got: %+v", i, item.expectedUser, user)
		}
	}
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUsersRepo(db)

	cases := []struct {
		expectedUser model.User
		expectedErr  error
		run          func(user model.User) (model.User, error)
	}{
		{
			expectedUser: model.User{ID: "1"},
			expectedErr:  nil,
			run: func(user model.User) (model.User, error) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"})
				rows.AddRow(user.ID, user.Username, user.Password)
				mock.
					ExpectQuery("SELECT id, username, password FROM user WHERE").
					WithArgs(user.ID).
					WillReturnRows(rows)
				return repo.GetUserByID(user.ID)
			},
		},
		{
			expectedUser: model.User{},
			expectedErr:  errors.New("bad query"),
			run: func(user model.User) (model.User, error) {
				mock.
					ExpectQuery("SELECT id, username, password FROM user WHERE").
					WithArgs(user.ID).
					WillReturnError(errors.New("bad query"))
				return repo.GetUserByID(user.ID)
			},
		},
		{
			expectedUser: model.User{},
			expectedErr:  customerr.UserNotFoundByID{UserID: "1"},
			run: func(user model.User) (model.User, error) {
				user.Username = "ivan"
				mock.
					ExpectQuery("SELECT id, username, password FROM user WHERE").
					WithArgs("1").
					WillReturnError(sql.ErrNoRows)
				return repo.GetUserByID("1")
			},
		},
	}

	for i, item := range cases {
		user, err := item.run(item.expectedUser)
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
		if !reflect.DeepEqual(item.expectedUser, user) {
			t.Errorf("[%d] expected user: %+v, got: %+v", i, item.expectedUser, user)
		}
	}
}
