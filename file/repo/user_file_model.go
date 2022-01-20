package repo

import (
	"encoding/json"
	"time"

	"github.com/changpro/disk-service/common/constants"
)

type UserFilePO struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string `bson:"user_id,omitempty" json:"userId,omitempty"`
	UniFileID string `bson:"uni_file_id,omitempty" json:"uniFileId,omitempty"`
	Name      string `bson:"name,omitempty" json:"name,omitempty"`
	FileMd5   string `bson:"file_md5,omitempty" json:"fileMd5,omitempty"`
	Path      string `bson:"path,omitempty" json:"path,omitempty"`
	IsDir     int32  `bson:"is_dir,omitempty" json:"isDir,omitempty"` // 2: file, 1: directory
	// IsDirEmpty int32     `bson:"is_dir_empty,omitempty" json:"isDirEmpty,omitempty"` // 2: not empty, 1: empty (case folder)
	// SubIDs   []string  `bson:"sub_ids,omitempty" json:"subIds,omitempty"`
	FileSize int64     `bson:"file_size,omitempty" json:"fileSize,omitempty"`
	FileType string    `bson:"file_type,omitempty" json:"fileType,omitempty"`
	IsHide   int32     `bson:"is_hide,omitempty" json:"isHide,omitempty"` // hide(1) or not(2)
	Status   int32     `bson:"status,omitempty" json:"status,omitempty"`  // enable(1), blacklist(2), recycle bin(3), deleted (4)
	CreateAt time.Time `bson:"create_at,omitempty" json:"-"`
	UpdateAt time.Time `bson:"update_at,omitempty" json:"-"`
}

func (p *UserFilePO) MarshalJSON() ([]byte, error) {
	type Alias UserFilePO
	return json.Marshal(&struct {
		*Alias
		CreateAt string `json:"createAt"`
		UpdateAt string `json:"updateAt"`
	}{
		Alias:    (*Alias)(p),
		CreateAt: p.CreateAt.Format(constants.StandardTimeFormat),
		UpdateAt: p.UpdateAt.Format(constants.StandardTimeFormat),
	})
}
