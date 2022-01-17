package interfaces

import (
	"context"

	cutil "github.com/changpro/disk-service/common/util"
	"github.com/changpro/disk-service/file/interfaces/assembler"
	"github.com/changpro/disk-service/file/service"
	filepb "github.com/changpro/disk-service/pbdeps/file"
)

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (s *server) UploadFile(ctx context.Context,
	req *filepb.UploadFileReq) (*filepb.UploadFileRsp, error) {
	rsp := &filepb.UploadFileRsp{}
	id, err := service.QuickUpload(ctx, req.UserId, req.FileName, req.FileMd5)
	if err != nil {
		return rsp, err
	}
	rsp.FileId = id
	return rsp, nil
}

func (s *server) GetDirsAndFiles(ctx context.Context,
	req *filepb.GetDirsAndFilesReq) (*filepb.GetDirsAndFilesRsp, error) {
	rsp := &filepb.GetDirsAndFilesRsp{}
	content, err := service.GetUserRoot(ctx, req.UserId)
	if err != nil {
		return rsp, err
	}
	for _, d := range content {
		rsp.Details = append(rsp.Details, cutil.AnyToStructpb(d))
	}
	return rsp, nil
}

func (s *server) GetFileDetail(ctx context.Context,
	req *filepb.GetFileDetailReq) (*filepb.GetFileDetailRsp, error) {
	rsp := &filepb.GetFileDetailRsp{}
	detail, err := service.GetFileDetail(ctx, req.UserId, req.FileId)
	if err != nil {
		return rsp, err
	}
	return assembler.AssembleFileDetail(detail), nil
}

func (s *server) MakeNewFolder(ctx context.Context,
	req *filepb.MakeNewFolderReq) (*filepb.MakeNewFolderRsp, error) {
	rsp := &filepb.MakeNewFolderRsp{}
	id, err := service.MakeNewFolder(ctx, req.UserId, req.DirName, req.Path)
	if err != nil {
		return rsp, err
	}
	rsp.Id = id
	return rsp, nil
}
