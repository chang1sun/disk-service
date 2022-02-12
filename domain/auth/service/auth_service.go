package service

import (
	"context"
	"strconv"
	"time"

	"github.com/changpro/disk-service/domain/auth/repo"
	"github.com/changpro/disk-service/infra/config"
	"github.com/changpro/disk-service/infra/constants"
	"github.com/changpro/disk-service/infra/errcode"
	"github.com/changpro/disk-service/infra/util"
	"github.com/golang-jwt/jwt"
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

func SignIn(ctx context.Context, userID string, password string) (string, error) {
	// add salt and calculate sha
	pwMask := util.GetStringWithSalt(password)
	user, err := repo.GetUserDao().SignIn(ctx, userID, pwMask)
	if err != nil {
		return "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if user == nil {
		return "", status.Error(errcode.NoSuchUserCode, errcode.NoSuchUserMsg)
	}
	now := time.Now()
	jwtId := userID + strconv.FormatInt(now.Unix(), 10)
	// set claims and sign
	claims := util.Claim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			Audience:  userID,
			ExpiresAt: now.Add(7 * 24 * time.Hour).Unix(),
			Id:        jwtId,
			IssuedAt:  now.Unix(),
			Issuer:    "easydisk",
			Subject:   userID,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(config.GetConfig().AuthKey))
	if err != nil {
		return "", status.Errorf(errcode.JWTParseErrCode, errcode.JWTParseErrMsg, err)
	}
	return token, nil
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
	if err := isAuthMatch(ctx, user, dto); err != nil {
		return err
	}
	if err := repo.GetUserDao().UpdatePassword(ctx, dto.UserID, dto.NewPw); err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func isAuthMatch(ctx context.Context, user *repo.UserPO, dto *repo.ModifyPwDTO) error {
	if !(user.AuthEmail == dto.AuthEmail || user.UserPW == dto.OldPw) {
		return status.Error(errcode.AuthMatchFailCode, errcode.AuthMatchFailMsg)
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
