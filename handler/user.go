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
	x           = "micro-starter"
	defaultCost = 10
)

var (
	alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

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
func (u *Users) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {

	if len(req.Password) < 8 {
		return errors.InternalServerError("users.Create.Check", "Password is less than 8 characters")
	}
	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Password), defaultCost)
	if err != nil {
		return errors.InternalServerError("user.Create", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	return u.repo.Create(&pb.User{
		Id:       req.Id,
		Username: strings.ToLower(req.Username),
		Email:    strings.ToLower(req.Email),
	}, salt, pp)
}

//Read a User info from repo
func (u *Users) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	user, err := u.repo.Read(req.Id)
	if err != nil {
		return err
	}
	rsp.User = user
	return nil
}

//Update user information
func (u *Users) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	return u.repo.Update(&pb.User{
		Id:       req.Id,
		Username: strings.ToLower(req.Username),
		Email:    strings.ToLower(req.Email),
	})
}

//Delete user information
func (u *Users) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	return u.repo.Delete(req.Id)
}

//Search user information
func (u *Users) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	users, err := u.repo.Search(req.Username, req.Email, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Users = users
	return nil
}

func random(i int) string {
	bytes := make([]byte, i)

	fix := byte(len(alphanum))
	_, _ = rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%fix]
	}
	return string(bytes)
}
