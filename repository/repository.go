package repository

import (
	"errors"
	"time"

	user "github.com/micro-community/micro-users/proto"
	"github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/store"
)

type pw struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type Table struct {
	users     model.Table
	sessions  model.Table
	passwords model.Table
}

func New() *Table {
	nameIndex := model.ByEquality("name")
	nameIndex.Unique = true
	emailIndex := model.ByEquality("email")
	emailIndex.Unique = true

	return &Table{
		users:     model.NewTable(store.DefaultStore, "users", model.Indexes(nameIndex, emailIndex), nil),
		sessions:  model.NewTable(store.DefaultStore, "sessions", nil, nil),
		passwords: model.NewTable(store.DefaultStore, "passwords", nil, nil),
	}
}

func (table *Table) CreateSession(sess *user.Session) error {
	if sess.Created == 0 {
		sess.Created = time.Now().Unix()
	}

	if sess.Expires == 0 {
		sess.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	return table.sessions.Save(sess)
}

func (table *Table) DeleteSession(id string) error {
	return table.sessions.Delete(model.Equals("id", id))
}

func (table *Table) ReadSession(id string) (*user.Session, error) {
	sess := &user.Session{}
	// @todo there should be a Read in the model to get rid of this pattern
	return sess, table.sessions.Read(model.Equals("id", id), &sess)
}

func (table *Table) Create(user *user.User, salt string, password string) error {
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()
	err := table.users.Save(user)
	if err != nil {
		return err
	}
	return table.passwords.Save(pw{
		ID:       user.Id,
		Password: password,
		Salt:     salt,
	})
}

func (table *Table) Delete(id string) error {
	return table.users.Delete(model.Equals("id", id))
}

func (table *Table) Update(user *user.User) error {
	user.Updated = time.Now().Unix()
	return table.users.Save(user)
}

func (table *Table) Read(id string) (*user.User, error) {
	user := &user.User{}
	return user, table.users.Read(model.Equals("id", id), user)
}

func (table *Table) Search(username, email string, limit, offset int64) ([]*user.User, error) {
	var query model.Query
	if len(username) > 0 {
		query = model.Equals("name", username)
	} else if len(email) > 0 {
		query = model.Equals("email", email)
	} else {
		return nil, errors.New("username and email cannot be blank")
	}

	users := []*user.User{}
	return users, table.users.List(query, &users)
}

func (table *Table) UpdatePassword(id string, salt string, password string) error {
	return table.passwords.Save(pw{
		ID:       id,
		Password: password,
		Salt:     salt,
	})
}

func (table *Table) SaltAndPassword(username, email string) (string, string, error) {
	var query model.Query
	if len(username) > 0 {
		query = model.Equals("name", username)
	} else if len(email) > 0 {
		query = model.Equals("email", email)
	} else {
		return "", "", errors.New("username and email cannot be blank")
	}

	user := &user.User{}
	err := table.users.Read(query, &user)
	if err != nil {
		return "", "", err
	}

	query = model.Equals("id", user.Id)
	query.Order.Type = model.OrderTypeUnordered

	password := &pw{}
	err = table.passwords.Read(query, password)
	if err != nil {
		return "", "", err
	}
	return password.Salt, password.Password, nil
}
