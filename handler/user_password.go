package handler

import (
	"encoding/base64"

	pb "github.com/micro-community/micro-users/proto"
	"github.com/micro/micro/v3/service/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

//UpdatePassword user information
func (s *Users) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest, rsp *pb.UpdatePasswordResponse) error {
	usr, err := s.repo.Read(req.UserId)
	if err != nil {
		return errors.InternalServerError("micro-community.srv.user.updatepassword", err.Error())
	}

	salt, hashed, err := s.repo.SaltAndPassword(usr.Username, usr.Email)
	if err != nil {
		return errors.InternalServerError("micro-community.srv.user.updatepassword", err.Error())
	}

	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("micro-community.srv.user.updatepassword", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.OldPassword)); err != nil {
		return errors.Unauthorized("micro-community.srv.user.updatepassword", err.Error())
	}

	salt = random(16) //use different salt
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.NewPassword), defaultCost)
	if err != nil {
		return errors.InternalServerError("micro-community.srv.user.updatepassword", err.Error())
	}
	encodedPassword := base64.StdEncoding.EncodeToString(h)

	if err := s.repo.UpdatePassword(req.UserId, salt, encodedPassword); err != nil {
		return errors.InternalServerError("micro-community.srv.user.updatepassword", err.Error())
	}
	return nil
}
