package repository

import (
	"errors"

	user "github.com/micro-community/micro-users/proto"
	"github.com/micro/dev/model"
)

//UpdatePassword of user
func (repo *Repos) UpdatePassword(id string, salt string, password string) error {
	return repo.passwords.Save(passwd{
		ID:       id,
		Password: password,
		Salt:     salt,
	})
}

//SaltAndPassword for user
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
