package assembler

import (
	"encoding/json"

	"github.com/changpro/disk-service/domain/file/repo"
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
			DocName:    record.DocName,
			Message:    record.Message,
			CreateTime: record.CreateTime.Format(constants.StandardTimeFormat),
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
		ExpireHour: detail.ExpireHours,
		CreateTime: detail.CreateTime,
		ViewNum:    detail.ViewNum,
		SaveNum:    detail.SaveNum,
		Size:       detail.DocSize,
	}
}
