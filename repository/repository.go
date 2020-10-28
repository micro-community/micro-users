package repository

import (
	"errors"
	"time"

	user "github.com/micro-community/micro-users/proto"
	"github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/store"
)

//Repos hold tables
type Repos struct {
	users     model.Table
	sessions  model.Table
	passwords model.Table

	nameIndex  model.Index
	emailIndex model.Index
	idIndex    model.Index
}

//New return a repo object
func New() *Repos {
	nameIndex := model.ByEquality("username")
	nameIndex.Unique = true
	nameIndex.Order.Type = model.OrderTypeUnordered

	emailIndex := model.ByEquality("email")
	emailIndex.Unique = true
	emailIndex.Order.Type = model.OrderTypeUnordered

	idIndex := model.ByEquality("id")
	idIndex.Order.Type = model.OrderTypeUnordered

	return &Repos{
		users:      model.NewTable(store.DefaultStore, "users", model.Indexes(nameIndex, emailIndex), nil),
		sessions:   model.NewTable(store.DefaultStore, "sessions", nil, nil),
		passwords:  model.NewTable(store.DefaultStore, "passwords", nil, nil),
		nameIndex:  nameIndex,
		emailIndex: emailIndex,
		idIndex:    idIndex,
	}
}

//Create a user,and save its password
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

//Delete a user by id
func (repo *Repos) Delete(id string) error {
	return repo.users.Delete(model.Equals("id", id))
}

//Update a user by model
func (repo *Repos) Update(user *user.User) error {
	user.Updated = time.Now().Unix()
	return repo.users.Save(user)
}

func (repo *Repos) Read(id string) (*user.User, error) {
	user := &user.User{}
	return user, repo.users.Read(model.Equals("id", id), user)
}

//Search user
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
