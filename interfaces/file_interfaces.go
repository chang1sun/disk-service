package interfaces

import (
	"context"
	"log"

	"github.com/changpro/disk-service/application"
	"github.com/changpro/disk-service/domain/file/service"
	"github.com/changpro/disk-service/infra/util"
	"github.com/changpro/disk-service/interfaces/assembler"
	"github.com/changpro/disk-service/stub"
	"google.golang.org/protobuf/types/known/emptypb"
)

type fileServerImpl struct {
}

func NewFileServer() *fileServerImpl {
	return &fileServerImpl{}
}

func (s *fileServerImpl) UploadFile(ctx context.Context,
	req *stub.UploadFileReq) (*stub.UploadFileRsp, error) {
	rsp := &stub.UploadFileRsp{}
	id, err := application.QuickUpload(ctx, req.UserId, req.FileName, req.FileMd5)
	if err != nil {
		util.LogErr(err, "UploadFile")
		return rsp, err
	}
	rsp.FileId = id
	return rsp, nil
}

func (s *fileServerImpl) GetDirsAndFiles(ctx context.Context,
	req *stub.GetDirsAndFilesReq) (*stub.GetDirsAndFilesRsp, error) {
	rsp := &stub.GetDirsAndFilesRsp{}
	log.Println(req.UserId, req.Path, req.ShowHide)
	content, err := service.GetDirByPath(ctx, req.UserId, req.Path, req.ShowHide)
	if err != nil {
		util.LogErr(err, "GetDirsAndFiles")
		return rsp, err
	}
	detail, err := assembler.AssembleDirsAndFilesList(content)
	if err != nil {
		util.LogErr(err, "GetDirsAndFiles")
		return rsp, err
	}
	rsp.Details = detail
	return rsp, nil
}

func (s *fileServerImpl) GetFileDetail(ctx context.Context,
	req *stub.GetFileDetailReq) (*stub.GetFileDetailRsp, error) {
	rsp := &stub.GetFileDetailRsp{}
	detail, err := service.GetFileDetail(ctx, req.UserId, req.FileId)
	if err != nil {
		util.LogErr(err, "GetFileDetail")
		return rsp, err
	}
	return assembler.AssembleFileDetail(detail), nil
}

func (s *fileServerImpl) MakeNewFolder(ctx context.Context,
	req *stub.MakeNewFolderReq) (*stub.MakeNewFolderRsp, error) {
	rsp := &stub.MakeNewFolderRsp{}
	id, err := service.MakeNewFolder(ctx, req.UserId, req.DirName, req.Path, req.Overwrite)
	if err != nil {
		util.LogErr(err, "MakeNewFolder")
		return rsp, err
	}
	rsp.Id = id
	return rsp, nil
}

func (s *fileServerImpl) Rename(ctx context.Context,
	req *stub.RenameReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.Rename(ctx, req.Id, req.NewName, req.Overwrite); err != nil {
		util.LogErr(err, "Rename")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) MoveToRecycleBin(ctx context.Context,
	req *stub.MoveToRecycleBinReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.MoveToRecycleBin(ctx, req.Id); err != nil {
		util.LogErr(err, "MoveToRecycleBin")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) SoftDelete(ctx context.Context,
	req *stub.SoftDeleteReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.SoftDelete(ctx, req.Id); err != nil {
		util.LogErr(err, "SoftDelete")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) HardDelete(ctx context.Context,
	req *stub.HardDeleteReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.HardDelete(ctx, req.Id); err != nil {
		util.LogErr(err, "HardDelete")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) CopyToPath(ctx context.Context,
	req *stub.CopyToPathReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.CopyToPath(ctx, req.Ids, req.Path, req.Overwrite)
	if err != nil {
		util.LogErr(err, "CopyToPath")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) MoveToPath(ctx context.Context,
	req *stub.MoveToPathReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.MoveToPath(ctx, req.Ids, req.Path, req.Overwrite)
	if err != nil {
		util.LogErr(err, "MoveToPath")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) CreateShare(ctx context.Context,
	req *stub.CreateShareReq) (*stub.CreateShareRsp, error) {
	rsp := &stub.CreateShareRsp{}
	token, err := service.CreateShare(ctx, assembler.AssembleCreateShareDTO(req))
	if err != nil {
		util.LogErr(err, "CreateShare")
		return rsp, err
	}
	rsp.Token = token
	return rsp, nil
}

func (s *fileServerImpl) RetrieveShareToPath(ctx context.Context,
	req *stub.RetrieveShareToPathReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := application.RetrieveShareFromToken(ctx, req.UserId, req.Token, req.Path)
	if err != nil {
		util.LogErr(err, "RetrieveShareToPath")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) GetShareRecords(ctx context.Context,
	req *stub.GetShareRecordsReq) (*stub.GetShareRecordsRsp, error) {
	rsp := &stub.GetShareRecordsRsp{}
	list, count, err := service.GetShareRecordList(ctx, req.UserId, req.Start, req.Limit)
	if err != nil {
		util.LogErr(err, "GetShareRecords")
		return rsp, err
	}
	rsp.Count = count
	rsp.List = assembler.AssembleShareRecordList(list)
	return rsp, nil
}

func (s *fileServerImpl) GetShareDetail(ctx context.Context,
	req *stub.GetShareDetailReq) (*stub.GetShareDetailRsp, error) {
	rsp := &stub.GetShareDetailRsp{}
	detail, err := service.GetShareDetail(ctx, req.Token)
	if err != nil {
		util.LogErr(err, "GetShareDetail")
		return rsp, err
	}
	return assembler.AssembleShareDetail(detail), nil
}

func (s *fileServerImpl) GetRecycleBinList(ctx context.Context,
	req *stub.GetRecycleBinListReq) (*stub.GetRecycleBinListRsp, error) {
	rsp := &stub.GetRecycleBinListRsp{}
	// TODO
	return rsp, nil
}