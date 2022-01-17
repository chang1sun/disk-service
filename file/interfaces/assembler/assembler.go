package assembler

import (
	"github.com/changpro/disk-service/common/constants"
	"github.com/changpro/disk-service/file/repo"
	filepb "github.com/changpro/disk-service/pbdeps/file"
)

func AssembleFileDetail(po *repo.UserFilePO) *filepb.GetFileDetailRsp {
	return &filepb.GetFileDetailRsp{
		Name:       po.FileName,
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
