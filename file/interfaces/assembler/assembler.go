package assembler

import (
	"encoding/json"

	"github.com/changpro/disk-service/common/constants"
	"github.com/changpro/disk-service/common/errcode"
	"github.com/changpro/disk-service/file/repo"
	filepb "github.com/changpro/disk-service/pbdeps/file"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func AssembleFileDetail(po *repo.UserFilePO) *filepb.GetFileDetailRsp {
	return &filepb.GetFileDetailRsp{
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
