package interfaces

import (
	"context"
	"log"

	"github.com/changpro/disk-service/domain/auth/service"
	"github.com/changpro/disk-service/interfaces/assembler"
	"github.com/changpro/disk-service/stub"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type authServerImpl struct {
}

func NewAuthServer() *authServerImpl {
	return &authServerImpl{}
}

func (*authServerImpl) RegisterNewUser(ctx context.Context,
	req *stub.RegisterNewUserReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.RegisterNewUser(ctx, assembler.AssembleUserPO(req))
	if err != nil {
		log.Println("RegisterNewUser failed, err msg: ", err)
		return rsp, err
	}
	return rsp, nil
}

func (*authServerImpl) SignIn(ctx context.Context,
	req *stub.SignInReq) (*stub.SignInRsp, error) {
	rsp := &stub.SignInRsp{}
	token, err := service.SignIn(ctx, req.UserId, req.Pw)
	if err != nil {
		log.Println("SignIn failed, err msg: ", err)
		return rsp, err
	}
	rsp.Token = token
	return rsp, nil
}

func (*authServerImpl) GetUserProfile(ctx context.Context,
	req *stub.GetUserProfileReq) (*stub.GetUserProfileRsp, error) {
	rsp := &stub.GetUserProfileRsp{}
	profile, err := service.GetUserProfile(ctx, req.UserId)
	if err != nil {
		return rsp, err
	}
	return assembler.AssembleUserProfile(profile), nil
}

func (*authServerImpl) ModifyPassword(ctx context.Context,
	req *stub.ModifyPasswordReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.ModifyPassword(ctx, assembler.AssembleModifyPwDTO(req))
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

func (*authServerImpl) ModifyUserProfile(ctx context.Context,
	req *stub.ModifyUserProfileReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.ModifyUserProfile(ctx, assembler.AssembleModifyUserProfileDTO(req))
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

func (*authServerImpl) UpdateUserStorage(ctx context.Context,
	req *stub.UpdateUserStorageReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.UpdateUserStorage(ctx, assembler.AssembleUpdateUserAnalysisDTO(req))
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}
