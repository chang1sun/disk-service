package interfaces

import (
	"context"

	authpb "github.com/changpro/disk-service/auth/deps"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (*server) RegisterNewUser(ctx context.Context,
	req *authpb.RegisterNewUserReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (*server) SignIn(ctx context.Context,
	req *authpb.SignInReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
