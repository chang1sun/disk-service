package auth

import (
	"context"
	"log"

	"github.com/changpro/disk-service/file/config"
	authpb "github.com/changpro/disk-service/pbdeps/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var caller *AuthCaller

type AuthCaller struct {
	client authpb.AuthServiceClient
}

func SetAuthCaller() {
	conn, err := grpc.Dial(config.GetConfig().AuthAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	caller = &AuthCaller{client: authpb.NewAuthServiceClient(conn)}
}

func GetAuthCaller() *AuthCaller {
	return caller
}

func (c *AuthCaller) UpdateUserStorage(ctx context.Context, userID string, fileNum int32, uploadFileNum int32, size int64) error {
	_, err := c.client.UpdateUserStorage(ctx, &authpb.UpdateUserStorageReq{
		UserId:        userID,
		FileNum:       fileNum,
		UploadFileNum: uploadFileNum,
		Size:          size,
	})
	if err != nil {
		return err
	}
	return nil
}
