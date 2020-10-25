package repository

import (
	"errors"
	"time"

	user "github.com/micro-community/micro-users/proto"
	"github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/store"
)

type passwd struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

//Repos hold tables
type Repos struct {
	users     model.Table
	sessions  model.Table
	passwords model.Table
}

//New return a repo object
func New() *Repos {
	nameIndex := model.ByEquality("name")
	nameIndex.Unique = true
	emailIndex := model.ByEquality("email")
	emailIndex.Unique = true

	return &Repos{
		users:     model.NewTable(store.DefaultStore, "users", model.Indexes(nameIndex, emailIndex), nil),
		sessions:  model.NewTable(store.DefaultStore, "sessions", nil, nil),
		passwords: model.NewTable(store.DefaultStore, "passwords", nil, nil),
	}
}

func (repo *Repos) CreateSession(session *user.Session) error {
	if session.Created == 0 {
		session.Created = time.Now().Unix()
	}

	if session.Expires == 0 {
		session.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	return repo.sessions.Save(session)
}

func (repo *Repos) DeleteSession(id string) error {
	return repo.sessions.Delete(model.Equals("id", id))
}

func (repo *Repos) ReadSession(id string) (*user.Session, error) {
	session := &user.Session{}
	// @todo there should be a Read in the model to get rid of this pattern
	return session, repo.sessions.Read(model.Equals("id", id), &session)
}

func (repo *Repos) Create(user *user.User, salt string, password string) error {
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()
	err := repo.users.Save(user)
	if err != nil {
		return err
	}
	return repo.passwords.Save(passwd{
		ID:       user.Id,
		Password: password,
		Salt:     salt,
	})
}

func (repo *Repos) Delete(id string) error {
	return repo.users.Delete(model.Equals("id", id))
}

func (repo *Repos) Update(user *user.User) error {
	user.Updated = time.Now().Unix()
	return repo.users.Save(user)
}

func (repo *Repos) Read(id string) (*user.User, error) {
	user := &user.User{}
	return user, repo.users.Read(model.Equals("id", id), user)
}

func (repo *Repos) Search(username, email string, limit, offset int64) ([]*user.User, error) {
	var query model.Query
	if len(username) > 0 {
		query = model.Equals("name", username)
	} else if len(email) > 0 {
		query = model.Equals("email", email)
	} else {
		return nil, errors.New("username and email cannot be blank")
	}

	users := []*user.User{}
	return users, repo.users.List(query, &users)
}

func (repo *Repos) UpdatePassword(id string, salt string, password string) error {
	return repo.passwords.Save(passwd{
		ID:       id,
		Password: password,
		Salt:     salt,
	})
}

func (repo *Repos) SaltAndPassword(username, email string) (string, string, error) {
	var query model.Query
	if len(username) > 0 {
		query = model.Equals("name", username)
	} else if len(email) > 0 {
		query = model.Equals("email", email)
	} else {
		return "", "", errors.New("username and email cannot be blank")
	}

	user := &user.User{}
	err := repo.users.Read(query, &user)
	if err != nil {
		return "", "", err
	}

	query = model.Equals("id", user.Id)
	query.Order.Type = model.OrderTypeUnordered

	password := &passwd{}
	err = repo.passwords.Read(query, password)
	if err != nil {
		return "", "", err
	}
	return password.Salt, password.Password, nil
}
