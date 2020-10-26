package repository

import (
	"time"

	user "github.com/micro-community/micro-users/proto"
	"github.com/micro/dev/model"
)

//CreateSession for user who login in successfully
func (repo *Repos) CreateSession(session *user.Session) error {
	if session.Created == 0 {
		session.Created = time.Now().Unix()
	}

	if session.Expires == 0 {
		session.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	return repo.sessions.Save(session)
}

//DeleteSession of user
func (repo *Repos) DeleteSession(id string) error {
	return repo.sessions.Delete(model.Equals("id", id))
}

//ReadSession by user id
func (repo *Repos) ReadSession(id string) (*user.Session, error) {
	session := &user.Session{}
	// @todo there should be a Read in the model to get rid of this pattern
	return session, repo.sessions.Read(model.Equals("id", id), &session)
}
