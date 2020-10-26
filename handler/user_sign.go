package handler

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/google/uuid"
	pb "github.com/micro-community/micro-users/proto"
	"github.com/micro/micro/v3/service/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

//Login user information
func (s *Users) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	username := strings.ToLower(req.Username)
	email := strings.ToLower(req.Email)

	salt, hashed, err := s.repo.SaltAndPassword(username, email)
	if err != nil {
		return err
	}
	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("micro-community.srv.user.Login", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.Password)); err != nil {
		return errors.Unauthorized("micro-community.srv.user.login", err.Error())
	}
	// save session
	session := &pb.Session{
		Id:       uuid.New().String(),
		Username: username,
		Created:  time.Now().Unix(),
		Expires:  time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	if err := s.repo.CreateSession(session); err != nil {
		return errors.InternalServerError("micro-community.srv.user.Login", err.Error())
	}
	rsp.Session = session
	return nil
}

//Logout user information
func (s *Users) Logout(ctx context.Context, req *pb.LogoutRequest, rsp *pb.LogoutResponse) error {
	return s.repo.DeleteSession(req.SessionId)
}
