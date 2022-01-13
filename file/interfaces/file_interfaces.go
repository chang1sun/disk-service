package interfaces

import (
	"context"

	filepb "github.com/changpro/disk-service/pbdeps/file"
)

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (s *server) UploadFile(ctx context.Context,
	req *filepb.UploadFileReq) (*filepb.UploadFileRsp, error) {
	return &filepb.UploadFileRsp{}, nil
}

func (s *server) GetDirsAndFiles(ctx context.Context,
	req *filepb.GetDirsAndFilesReq) (*filepb.GetDirsAndFilesRsp, error) {
	return &filepb.GetDirsAndFilesRsp{}, nil
}

func (s *server) GetFileDetail(ctx context.Context,
	req *filepb.GetFileDetailReq) (*filepb.GetFileDetailRsp, error) {
	return &filepb.GetFileDetailRsp{}, nil
}

func (s *server) DownloadFile(ctx context.Context,
	req *filepb.DownloadFileReq) (*filepb.DownloadFileRsp, error) {
	return &filepb.DownloadFileRsp{}, nil
}
