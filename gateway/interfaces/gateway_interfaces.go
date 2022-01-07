package interfaces

import (
	"context"

	gatewaypb "github.com/changpro/disk-service/gateway/deps"
)

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (s *server) RegisterNewUser(ctx context.Context,
	req *gatewaypb.RegisterNewUserReq) (*gatewaypb.CommonHttpRsp, error) {
	return &gatewaypb.CommonHttpRsp{}, nil
}

func (s *server) SignIn(ctx context.Context,
	req *gatewaypb.SignInReq) (*gatewaypb.CommonHttpRsp, error) {
	return &gatewaypb.CommonHttpRsp{}, nil
}

func (s *server) SetIcon(ctx context.Context,
	req *gatewaypb.SetIconReq) (*gatewaypb.CommonHttpRsp, error) {
	return &gatewaypb.CommonHttpRsp{}, nil
}

func (s *server) ModifyPassword(ctx context.Context,
	req *gatewaypb.ModifyPasswordReq) (*gatewaypb.CommonHttpRsp, error) {
	return &gatewaypb.CommonHttpRsp{}, nil
}

func (s *server) GetUserProfile(ctx context.Context,
	req *gatewaypb.GetUserProfileReq) (*gatewaypb.CommonHttpRsp, error) {
	return &gatewaypb.CommonHttpRsp{}, nil
}
