package interfaces

import (
	"context"

	authpb "github.com/changpro/disk-service/auth/deps"
)

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (*server) RegisterNewUser(ctx context.Context,
	req *authpb.RegisterNewUserReq)
