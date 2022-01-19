package interfaces

import (
	"context"

	cutil "github.com/changpro/disk-service/common/util"
	"github.com/changpro/disk-service/file/interfaces/assembler"
	"github.com/changpro/disk-service/file/service"
	filepb "github.com/changpro/disk-service/pbdeps/file"
	"google.golang.org/protobuf/types/known/emptypb"
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
		cutil.LogErr(err, "UploadFile")
		return rsp, err
	}
	rsp.FileId = id
	return rsp, nil
}

func (s *server) GetDirsAndFiles(ctx context.Context,
	req *filepb.GetDirsAndFilesReq) (*filepb.GetDirsAndFilesRsp, error) {
	rsp := &filepb.GetDirsAndFilesRsp{}
	content, err := service.GetDirByPath(ctx, req.UserId, req.Path, req.ShowHide)
	if err != nil {
		cutil.LogErr(err, "GetDirsAndFiles")
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
		cutil.LogErr(err, "GetFileDetail")
		return rsp, err
	}
	return assembler.AssembleFileDetail(detail), nil
}

func (s *server) MakeNewFolder(ctx context.Context,
	req *filepb.MakeNewFolderReq) (*filepb.MakeNewFolderRsp, error) {
	rsp := &filepb.MakeNewFolderRsp{}
	id, err := service.MakeNewFolder(ctx, req.UserId, req.DirName, req.Path, req.Overwrite)
	if err != nil {
		cutil.LogErr(err, "MakeNewFolder")
		return rsp, err
	}
	rsp.Id = id
	return rsp, nil
}

func (s *server) Rename(ctx context.Context,
	req *filepb.RenameReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.Rename(ctx, req.Id, req.NewName, req.Overwrite); err != nil {
		cutil.LogErr(err, "Rename")
		return rsp, err
	}
	return rsp, nil
}

func (s *server) MoveToRecycleBin(ctx context.Context,
	req *filepb.MoveToRecycleBinReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.MoveToRecycleBin(ctx, req.Id); err != nil {
		cutil.LogErr(err, "MoveToRecycleBin")
		return rsp, err
	}
	return rsp, nil
}

func (s *server) SoftDelete(ctx context.Context,
	req *filepb.SoftDeleteReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.SoftDelete(ctx, req.Id); err != nil {
		cutil.LogErr(err, "SoftDelete")
		return rsp, err
	}
	return rsp, nil
}

func (s *server) HardDelete(ctx context.Context,
	req *filepb.HardDeleteReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.HardDelete(ctx, req.Id); err != nil {
		cutil.LogErr(err, "HardDelete")
		return rsp, err
	}
	return rsp, nil
}

func (s *server) CopyToPath(ctx context.Context,
	req *filepb.CopyToPathReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.CopyToPath(ctx, req.Ids, req.Path, req.Overwrite)
	if err != nil {
		cutil.LogErr(err, "CopyToPath")
		return rsp, err
	}
	return rsp, nil
}

func (s *server) MoveToPath(ctx context.Context,
	req *filepb.MoveToPathReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.MoveToPath(ctx, req.Ids, req.Path, req.Overwrite)
	if err != nil {
		cutil.LogErr(err, "MoveToPath")
		return rsp, err
	}
	return rsp, nil
}

func (s *server) CreateShare(ctx context.Context,
	req *filepb.CreateShareReq) (*filepb.CreateShareRsp, error) {
	rsp := &filepb.CreateShareRsp{}
	// TODO
	return rsp, nil
}

func (s *server) RetrieveShareToPath(ctx context.Context,
	req *filepb.RetrieveShareToPathReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	// TODO
	return rsp, nil
}

func (s *server) GetRecycleBinList(ctx context.Context,
	req *filepb.GetRecycleBinListReq) (*filepb.GetRecycleBinListRsp, error) {
	rsp := &filepb.GetRecycleBinListRsp{}
	// TODO
	return rsp, nil
}

func (s *server) GetShareRecords(ctx context.Context,
	req *filepb.GetShareRecordsReq) (*filepb.GetShareRecordsRsp, error) {
	rsp := &filepb.GetShareRecordsRsp{}
	// TODO
	return rsp, nil
}
