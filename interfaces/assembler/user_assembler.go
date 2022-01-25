package assembler

import (
	"github.com/changpro/disk-service/domain/auth/repo"
	"github.com/changpro/disk-service/domain/auth/service"
	"github.com/changpro/disk-service/stub"
)

func AssembleUserPO(req *stub.RegisterNewUserReq) *repo.UserPO {
	return &repo.UserPO{
		UserID:    req.UserId,
		UserPW:    req.Pw,
		AuthEmail: req.AuthEmail,
	}
}

func AssembleUserProfile(profile *service.UserProfile) *stub.GetUserProfileRsp {
	return &stub.GetUserProfileRsp{
		Icon:          profile.Icon,
		AuthEmail:     profile.AuthEmail,
		RegisterTime:  profile.RegisterAt,
		FileNum:       profile.FileNum,
		FileUploadNum: profile.FileUploadNum,
		TotalSize:     profile.TotalSize,
		UsedSize:      profile.UsedSize,
	}
}

func AssembleModifyPwDTO(req *stub.ModifyPasswordReq) *repo.ModifyPwDTO {
	return &repo.ModifyPwDTO{
		UserID:    req.UserId,
		NewPw:     req.NewPw,
		AuthEmail: req.AuthEmail,
		OldPw:     req.OldPw,
	}
}

func AssembleModifyUserProfileDTO(req *stub.ModifyUserProfileReq) *repo.ModifyUserProfileDTO {
	return &repo.ModifyUserProfileDTO{
		UserID:    req.UserId,
		AuthEmail: req.AuthEmail,
		Icon:      req.Icon,
	}
}

func AssembleUpdateUserAnalysisDTO(req *stub.UpdateUserStorageReq) *repo.UpdateUserAnalysisDTO {
	return &repo.UpdateUserAnalysisDTO{
		UserID:        req.UserId,
		FileNum:       req.FileNum,
		Size:          req.Size,
		UploadFileNum: req.UploadFileNum,
	}
}
