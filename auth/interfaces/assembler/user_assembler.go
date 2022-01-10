package assembler

import (
	authpb "github.com/changpro/disk-service/auth/deps"
	"github.com/changpro/disk-service/auth/repo"
	"github.com/changpro/disk-service/auth/service"
)

func AssembleUserPO(req *authpb.RegisterNewUserReq) *repo.UserPO {
	return &repo.UserPO{
		UserID:    req.UserId,
		UserPW:    req.Pw,
		AuthEmail: req.AuthEmail,
	}
}

func AssembleUserProfile(profile *service.UserProfile) *authpb.GetUserProfileRsp {
	return &authpb.GetUserProfileRsp{
		Icon:          profile.Icon,
		AuthEmail:     profile.AuthEmail,
		RegisterTime:  profile.RegisterAt,
		FileNum:       profile.FileNum,
		FileUploadNum: profile.FileUploadNum,
		TotalSize:     profile.TotalSize,
		UsedSize:      profile.UsedSize,
	}
}

func AssembleModifyPwDTO(req *authpb.ModifyPasswordReq) *repo.ModifyPwDTO {
	return &repo.ModifyPwDTO{
		UserID:    req.UserId,
		NewPw:     req.NewPw,
		AuthEmail: req.AuthEmail,
		OldPw:     req.OldPw,
	}
}

func AssembleModifyUserProfileDTO(req *authpb.ModifyUserProfileReq) *repo.ModifyUserProfileDTO {
	return &repo.ModifyUserProfileDTO{
		UserID:    req.UserId,
		AuthEmail: req.AuthEmail,
		Icon:      req.Icon,
	}
}
