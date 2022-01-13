package repo

import "time"

type UserFilePO struct {
	ID         string       `bson:"_id,omitempty"`
	Name       string       `bson:"name,omitempty"`
	FileMd5    string       `bson:"file_md5,omitempty"`
	IsDir      int32        `bson:"is_dir,omitempty"`       // 0: file, 1: directory
	IsDirEmpty int32        `bson:"is_dir_empty,omitempty"` // 0: not empty, 1: empty
	FileSize   int64        `bson:"file_size,omitempty"`
	UploadAt   time.Time    `bson:"upload_at,omitempty"`
	UploadBy   string       `bson:"upload_by,omitempty"`
	AddAt      time.Time    `bson:"add_at,omitempty"`
	UpdateAt   time.Time    `bson:"update_at,omitempty"`
	Sub        []UserFilePO `bson:"sub,omitempty"`
}

type UserDirPO struct {
	UserID string     `bson:"user_id,omitempty"`
	Root   UserFilePO `bson:"root,omitempty"`
}
