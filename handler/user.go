package handler

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	pb "github.com/micro-community/micro-users/proto"
	repo "github.com/micro-community/micro-users/repository"
	"github.com/micro/micro/v3/service/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

const (
	x = "micro-starter"
)

var (
	alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func random(i int) string {
	bytes := make([]byte, i)

	fix := byte(len(alphanum))
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%fix]
	}
	return string(bytes)
	//return "ughwhy?!!!"
}

//Users handler
type Users struct {
	repo *repo.Repos
}

//NewUsers Return users handler
func NewUsers() *Users {
	return &Users{
		repo: repo.New(),
	}
}

//Create a user
func (s *Users) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Password), 10)
	if err != nil {
		return errors.InternalServerError("micro-community.srv.user.Create", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	req.User.Username = strings.ToLower(req.User.Username)
	req.User.Email = strings.ToLower(req.User.Email)
	return s.repo.Create(req.User, salt, pp)
}

//Read s User
func (s *Users) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	user, err := s.repo.Read(req.Id)
	if err != nil {
		return err
	}
	rsp.User = user
	return nil
}

//Update user information
func (s *Users) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	req.User.Username = strings.ToLower(req.User.Username)
	req.User.Email = strings.ToLower(req.User.Email)
	return s.repo.Update(req.User)
}

//Delete user information
func (s *Users) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	return s.repo.Delete(req.Id)
}

//Search user information
func (s *Users) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	users, err := s.repo.Search(req.Username, req.Email, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Users = users
	return nil
}
