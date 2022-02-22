package assembler

import (
	"encoding/json"
	"time"

	"github.com/changpro/disk-service/domain/file/repo"
	"github.com/changpro/disk-service/domain/file/service"
	"github.com/changpro/disk-service/infra/constants"
	"github.com/changpro/disk-service/infra/errcode"
	"github.com/changpro/disk-service/stub"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func AssembleFileDetail(po *repo.UserFilePO) *stub.GetFileDetailRsp {
	return &stub.GetFileDetailRsp{
		Name:       po.Name,
		Id:         po.ID,
		Size:       po.FileSize,
		Type:       po.FileType,
		Status:     po.Status,
		Md5:        po.FileMd5,
		Path:       po.Path,
		CreateTime: po.CreateAt.Format(constants.StandardTimeFormat),
		UpdateTime: po.UpdateAt.Format(constants.StandardTimeFormat),
	}
}

func AssembleDirsAndFilesList(pos []*repo.UserFilePO) ([]*structpb.Struct, error) {
	var list []*structpb.Struct
	for _, po := range pos {
		data, err := json.Marshal(po)
		if err != nil {
			return nil, status.Errorf(errcode.JsonMarshalErrCode, errcode.JsonMarshalErrMsg, err)
		}
		m := make(map[string]interface{})
		err = json.Unmarshal(data, &m)
		if err != nil {
			return nil, status.Errorf(errcode.JsonMarshalErrCode, errcode.JsonMarshalErrMsg, err)
		}
		s, err := structpb.NewStruct(m)
		if err != nil {
			return nil, status.Errorf(errcode.StructpbConvertErrCode, errcode.StructpbConvertErrMsg, err)
		}
		list = append(list, s)
	}
	return list, nil
}

func AssembleCreateShareDTO(req *stub.CreateShareReq) *repo.CreateShareDTO {
	return &repo.CreateShareDTO{
		UserID:     req.UserId,
		DocID:      req.DocId,
		ExpireHour: req.ExpireHour,
	}
}

func AssembleShareRecordList(records []*repo.ShareRecordPO) []*stub.ShareRecord {
	var list []*stub.ShareRecord
	for _, record := range records {
		list = append(list, &stub.ShareRecord{
			Id:         record.ID,
			DocId:      record.DocID,
			DocName:    record.DocName,
			CreateTime: record.CreateTime.UnixMilli(),
			ExpireTime: record.ExpireTime.UnixMilli(),
			Token:      record.Token,
			Type:       record.Type,
			Status:     record.Status,
		})
	}
	return list
}

func AssembleShareDetail(detail *repo.ShareDetailPO) *stub.GetShareDetailRsp {
	return &stub.GetShareDetailRsp{
		Uploader:   detail.Uploader,
		DocId:      detail.DocID,
		DocName:    detail.DocName,
		DocType:    detail.DocType,
		IsDir:      detail.IsDir,
		ExpireHour: detail.ExpireHours,
		CreateTime: detail.CreateTime,
		ViewNum:    detail.ViewNum,
		SaveNum:    detail.SaveNum,
		DocSize:    detail.DocSize,
		FileNum:    detail.FileNum,
		UniFileId:  detail.UniFileID,
	}
}

func AssembleRecycleDocList(list []*repo.RecycleFilePO) []*stub.RecycleDocInfo {
	var res []*stub.RecycleDocInfo
	for _, doc := range list {
		res = append(res, &stub.RecycleDocInfo{
			DocId:    doc.ID,
			DocName:  doc.Name,
			IsDir:    doc.IsDir,
			DeleteAt: doc.DeleteAt.Format("2006-01-02"),
			ExpireAt: doc.ExpireAt.Format("2006-01-02"),
		})
	}
	return res
}

func AssemblShareRecordQuery(req *stub.GetShareRecordsReq) *repo.RecordQuery {
	return &repo.RecordQuery{
		UserID:    req.UserId,
		Offset:    req.Offset,
		Limit:     req.Offset,
		Type:      req.Type,
		StartTime: time.UnixMilli(req.StartTime).Unix(),
		EndTime:   time.UnixMilli(req.EndTime).Unix(),
	}
}

func AssemblClassQuery(req *stub.GetClassifiedDocsReq) *repo.ClassifiedDocsQuery {
	return &repo.ClassifiedDocsQuery{
		UserID: req.UserId,
		Type:   req.Type,
		Offset: req.Offset,
		Limit:  req.Limit,
	}
}

func AssemblClassifiedDocList(list []*repo.UserFilePO) []*stub.ClassifiedDoc {
	var res []*stub.ClassifiedDoc
	for _, po := range list {
		res = append(res, &stub.ClassifiedDoc{
			DocId:    po.ID,
			DocName:  po.Name,
			DocSize:  po.FileSize,
			DocType:  po.FileType,
			DocPath:  po.Path,
			CreateAt: po.CreateAt.Format(constants.StandardTimeFormat),
			UpdateAt: po.UpdateAt.Format(constants.StandardTimeFormat),
		})
	}
	return res
}

func AssemblShareTree(root *service.ShareTreeNode) *stub.ShareFolderTreeNode {
	node := &stub.ShareFolderTreeNode{
		DocId:     root.DocID,
		DocName:   root.DocName,
		UniFileId: root.UniFileID,
		DocSize:   root.DocSize,
		IsDir:     root.IsDir,
	}
	var children []*stub.ShareFolderTreeNode
	for _, sub := range root.Children {
		children = append(children, AssemblShareTree(sub))
	}
	node.Children = children
	return node
}
