package handler

import (
	pb "github.com/micro-community/micro-users/proto"
	"golang.org/x/net/context"
)

//ReadSession user information
func (s *Users) ReadSession(ctx context.Context, req *pb.ReadSessionRequest, rsp *pb.ReadSessionResponse) error {
	session, err := s.repo.ReadSession(req.SessionId)
	if err != nil {
		return err
	}
	rsp.Session = session
	return nil
}
