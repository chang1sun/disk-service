package interfaces

import (
	"context"

	gatewaypb "github.com/changpro/disk-service/resource/stub"
)

type gatewayServerImpl struct {
}

func NewServer() *gatewayServerImpl {
	return &gatewayServerImpl{}
}

func (s *gatewayServerImpl) RegisterNewUser(ctx context.Context,
	req *gatewaypb.RegisterNewUserReq) (*gatewaypb.CommonHttpRsp, error) {
	return &gatewaypb.CommonHttpRsp{}, nil
}

func (s *gatewayServerImpl) SignIn(ctx context.Context,
	req *gatewaypb.SignInReq) (*gatewaypb.CommonHttpRsp, error) {
	return &gatewaypb.CommonHttpRsp{}, nil
}

func (s *gatewayServerImpl) SetIcon(ctx context.Context,
	req *gatewaypb.SetIconReq) (*gatewaypb.CommonHttpRsp, error) {
	return &gatewaypb.CommonHttpRsp{}, nil
}

func (s *gatewayServerImpl) ModifyPassword(ctx context.Context,
	req *gatewaypb.ModifyPasswordReq) (*gatewaypb.CommonHttpRsp, error) {
	return &gatewaypb.CommonHttpRsp{}, nil
}

func (s *gatewayServerImpl) GetUserProfile(ctx context.Context,
	req *gatewaypb.GetUserProfileReq) (*gatewaypb.CommonHttpRsp, error) {
	return &gatewaypb.CommonHttpRsp{}, nil
}
