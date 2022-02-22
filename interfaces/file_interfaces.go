package interfaces

import (
	"context"

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
	return rsp, nil
}

func (s *fileServerImpl) GetDirsAndFiles(ctx context.Context,
	req *stub.GetDirsAndFilesReq) (*stub.GetDirsAndFilesRsp, error) {
	rsp := &stub.GetDirsAndFilesRsp{}
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

func (s *fileServerImpl) SetHiddenDoc(ctx context.Context,
	req *stub.SetHiddenDocReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.SetHiddenDoc(ctx, req.Ids, req.HideStatus)
	if err != nil {
		util.LogErr(err, "SetHiddenDoc")
		return rsp, err
	}
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
	if err := service.MoveToRecycleBin(ctx, req.UserId, req.Ids); err != nil {
		util.LogErr(err, "MoveToRecycleBin")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) RecoverDocs(ctx context.Context,
	req *stub.RecoverDocsReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.RecoverDocs(ctx, req.UserId, req.Ids); err != nil {
		util.LogErr(err, "RecoverDocs")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) SoftDelete(ctx context.Context,
	req *stub.SoftDeleteReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.SoftDelete(ctx, req.UserId, req.Ids); err != nil {
		util.LogErr(err, "SoftDelete")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) HardDelete(ctx context.Context,
	req *stub.HardDeleteReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.HardDelete(ctx, req.Ids); err != nil {
		util.LogErr(err, "HardDelete")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) CopyToPath(ctx context.Context,
	req *stub.CopyToPathReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.CopyToPath(ctx, req.UserId, req.Ids, req.Path, req.Overwrite)
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
	token, password, err := service.CreateShare(ctx, assembler.AssembleCreateShareDTO(req))
	if err != nil {
		util.LogErr(err, "CreateShare")
		return rsp, err
	}
	rsp.Token = token
	rsp.Password = password
	return rsp, nil
}

func (s *fileServerImpl) RetrieveShareToPath(ctx context.Context,
	req *stub.RetrieveShareToPathReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	err := service.RetrieveShareFromToken(ctx, req.UserId, req.Token, req.Path)
	if err != nil {
		util.LogErr(err, "RetrieveShareToPath")
		return rsp, err
	}
	return rsp, nil
}

func (s *fileServerImpl) GetShareRecords(ctx context.Context,
	req *stub.GetShareRecordsReq) (*stub.GetShareRecordsRsp, error) {
	rsp := &stub.GetShareRecordsRsp{}
	list, count, err := service.GetShareRecordList(ctx, assembler.AssemblShareRecordQuery(req))
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
	detail, err := service.GetShareDetail(ctx, req.Token, req.Password)
	if err != nil {
		util.LogErr(err, "GetShareDetail")
		return rsp, err
	}
	return assembler.AssembleShareDetail(detail), nil
}

func (s *fileServerImpl) GetRecycleBinList(ctx context.Context,
	req *stub.GetRecycleBinListReq) (*stub.GetRecycleBinListRsp, error) {
	rsp := &stub.GetRecycleBinListRsp{}
	list, err := service.GetRecycleBinList(ctx, req.UserId, req.Offset, req.Limit)
	if err != nil {
		util.LogErr(err, "GetRecycleBinList")
		return rsp, err
	}
	rsp.List = assembler.AssembleRecycleDocList(list)
	return rsp, nil
}

func (s *fileServerImpl) GetClassifiedDocs(ctx context.Context,
	req *stub.GetClassifiedDocsReq) (*stub.GetClassifiedDocsRsp, error) {
	rsp := &stub.GetClassifiedDocsRsp{}
	list, err := service.GetClassifiedDocs(ctx, assembler.AssemblClassQuery(req))
	if err != nil {
		util.LogErr(err, "GetClassifiedDocs")
		return rsp, err
	}
	rsp.List = assembler.AssemblClassifiedDocList(list)
	return rsp, nil
}

func (s *fileServerImpl) GetShareGlimpse(ctx context.Context,
	req *stub.GetShareGlimpseReq) (*stub.GetShareGlimpseRsp, error) {
	rsp := &stub.GetShareGlimpseRsp{}
	uploader, docName, err := service.GetShareGlimpse(ctx, req.Token)
	if err != nil {
		util.LogErr(err, "GetShareGlimpse")
		return rsp, err
	}
	rsp.Uploader = uploader
	rsp.DocName = docName
	return rsp, nil
}

func (s *fileServerImpl) GetShareFolderTree(ctx context.Context,
	req *stub.GetShareFolderTreeReq) (*stub.GetShareFolderTreeRsp, error) {
	rsp := &stub.GetShareFolderTreeRsp{}
	root, err := service.GetShareFolderTree(ctx, req.Uploader, req.DocId)
	if err != nil {
		util.LogErr(err, "GetShareFolderTree")
		return rsp, err
	}
	rsp.Root = assembler.AssemblShareTree(root)
	return rsp, nil
}

func (s *fileServerImpl) DeleteShare(ctx context.Context,
	req *stub.DeleteShareReq) (*emptypb.Empty, error) {
	rsp := &emptypb.Empty{}
	if err := service.DeleteShare(ctx, req.Token); err != nil {
		util.LogErr(err, "DeleteShare")
		return rsp, err
	}
	return rsp, nil
}
