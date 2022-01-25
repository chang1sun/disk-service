package service

import (
	"context"

	"github.com/changpro/disk-service/auth/repo"
	"github.com/changpro/disk-service/auth/util"
	"github.com/changpro/disk-service/common/constants"
	"github.com/changpro/disk-service/common/errcode"
	"google.golang.org/grpc/status"
)

type UserProfile struct {
	Icon          string
	AuthEmail     string
	RegisterAt    string
	FileNum       int32
	FileUploadNum int32
	TotalSize     int64
	UsedSize      int64
}

func RegisterNewUser(ctx context.Context, userPO *repo.UserPO) error {
	// check repeat
	user, err := repo.GetUserDao().QueryUserByID(ctx, userPO.UserID)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if user != nil {
		return status.Error(errcode.DetectRepeatedUserIDCode, errcode.DetectRepeatedUserIDMsg)
	}

	// add salt and calculate sha
	pwMask := util.GetStringWithSalt(userPO.UserPW)
	userPO.UserPW = pwMask

	// insert a record
	err = repo.GetUserDao().RegisterNewUser(ctx, userPO)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func SignIn(ctx context.Context, userID string, password string) error {
	// add salt and calculate sha
	pwMask := util.GetStringWithSalt(password)
	user, err := repo.GetUserDao().SignIn(ctx, userID, pwMask)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if user == nil {
		return status.Error(errcode.NoSuchUserCode, errcode.NoSuchUserMsg)
	}
	return nil
}

func GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	user, err := repo.GetUserDao().QueryUserByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	analysis, err := repo.GetUserAnalysisDao().QueryUserAnalysisByUserID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return assembleUserProfile(user, analysis), nil
}

func assembleUserProfile(user *repo.UserPO, analysis *repo.UserAnalysisPO) *UserProfile {
	return &UserProfile{
		Icon:          user.UserIcon,
		AuthEmail:     user.AuthEmail,
		RegisterAt:    user.CreateTime.Format(constants.StandardTimeFormat),
		TotalSize:     analysis.TotalSize,
		UsedSize:      analysis.UsedSize,
		FileNum:       analysis.FileNum,
		FileUploadNum: analysis.FileNum,
	}
}

func ModifyPassword(ctx context.Context, dto *repo.ModifyPwDTO) error {
	dto.NewPw = util.GetStringWithSalt(dto.NewPw)
	dto.OldPw = util.GetStringWithSalt(dto.OldPw)
	user, err := repo.GetUserDao().QueryUserByID(ctx, dto.UserID)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if user == nil {
		return status.Error(errcode.NoSuchUserCode, errcode.NoSuchUserMsg)
	}
	if err := repo.GetUserDao().UpdatePassword(ctx, dto); err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func ModifyUserProfile(ctx context.Context, dto *repo.ModifyUserProfileDTO) error {
	if err := repo.GetUserDao().UpdateUserProfile(ctx, dto); err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func UpdateUserStorage(ctx context.Context, dto *repo.UpdateUserAnalysisDTO) error {
	// cal left vol is enough or not
	ana, err := repo.GetUserAnalysisDao().QueryUserAnalysisByUserID(ctx, dto.UserID)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if ana == nil {
		return status.Error(errcode.NoSuchUserCode, errcode.NoSuchUserMsg)
	}
	if ana.TotalSize-ana.UsedSize < dto.Size {
		return status.Error(errcode.NoEnoughVolCode, errcode.NoEnoughVolMsg)
	}
	// do update
	if err := repo.GetUserAnalysisDao().UpdateUserStorage(ctx, dto); err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}
