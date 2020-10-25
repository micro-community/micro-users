package handler

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

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

//Create user
func (s *Users) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Password), 10)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.Create", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	req.User.Username = strings.ToLower(req.User.Username)
	req.User.Email = strings.ToLower(req.User.Email)
	return s.repo.Create(req.User, salt, pp)
}

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

func (s *Users) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	return s.repo.Delete(req.Id)
}

func (s *Users) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	users, err := s.repo.Search(req.Username, req.Email, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Users = users
	return nil
}

func (s *Users) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest, rsp *pb.UpdatePasswordResponse) error {
	usr, err := s.repo.Read(req.UserId)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}

	salt, hashed, err := s.repo.SaltAndPassword(usr.Username, usr.Email)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}

	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.OldPassword)); err != nil {
		return errors.Unauthorized("go.micro.srv.user.updatepassword", err.Error())
	}

	salt = random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.NewPassword), 10)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	if err := s.repo.UpdatePassword(req.UserId, salt, pp); err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}
	return nil
}

func (s *Users) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	username := strings.ToLower(req.Username)
	email := strings.ToLower(req.Email)

	salt, hashed, err := s.repo.SaltAndPassword(username, email)
	if err != nil {
		return err
	}
	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.Login", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.Password)); err != nil {
		return errors.Unauthorized("go.micro.srv.user.login", err.Error())
	}
	// save session
	session := &pb.Session{
		Id:       random(128),
		Username: username,
		Created:  time.Now().Unix(),
		Expires:  time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	if err := s.repo.CreateSession(session); err != nil {
		return errors.InternalServerError("go.micro.srv.user.Login", err.Error())
	}
	rsp.Session = session
	return nil
}

func (s *Users) Logout(ctx context.Context, req *pb.LogoutRequest, rsp *pb.LogoutResponse) error {
	return s.repo.DeleteSession(req.SessionId)
}

func (s *Users) ReadSession(ctx context.Context, req *pb.ReadSessionRequest, rsp *pb.ReadSessionResponse) error {
	session, err := s.repo.ReadSession(req.SessionId)
	if err != nil {
		return err
	}
	rsp.Session = session
	return nil
}
