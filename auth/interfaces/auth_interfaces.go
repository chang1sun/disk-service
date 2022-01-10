package interfaces

import (
	"context"
	"log"

	authpb "github.com/changpro/disk-service/auth/deps"
	"github.com/changpro/disk-service/auth/interfaces/assembler"
	"github.com/changpro/disk-service/auth/service"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (*server) RegisterNewUser(ctx context.Context,
	req *authpb.RegisterNewUserReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.RegisterNewUser(ctx, assembler.AssembleUserPO(req))
	if err != nil {
		log.Println("RegisterNewUser failed, err msg: ", err)
		return rsp, err
	}
	return rsp, nil
}

func (*server) SignIn(ctx context.Context,
	req *authpb.SignInReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.SignIn(ctx, req.UserId, req.Pw)
	if err != nil {
		log.Println("SignIn failed, err msg: ", err)
		return rsp, err
	}
	return rsp, nil
}

func (*server) GetUserProfile(ctx context.Context,
	req *authpb.GetUserProfileReq) (*authpb.GetUserProfileRsp, error) {
	rsp := &authpb.GetUserProfileRsp{}
	profile, err := service.GetUserProfile(ctx, req.UserId)
	if err != nil {
		return rsp, err
	}
	return assembler.AssembleUserProfile(profile), nil
}

func (*server) ModifyPassword(ctx context.Context,
	req *authpb.ModifyPasswordReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.ModifyPassword(ctx, assembler.AssembleModifyPwDTO(req))
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

func (*server) ModifyUserProfile(ctx context.Context,
	req *authpb.ModifyUserProfileReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.ModifyUserProfile(ctx, assembler.AssembleModifyUserProfileDTO(req))
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}
