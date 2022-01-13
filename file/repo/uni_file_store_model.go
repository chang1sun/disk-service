package repo

import "time"

type UniFileMetaPO struct {
	Md5      string    `bson:"md5,omitempty"`
	Size     int64     `bson:"size,omitempty"`
	Type     string    `bson:"type,omitempty"`
	UploadAt time.Time `bson:"upload_at,omitempty"`
	UploadBy string    `bson:"upload_by,omitempty"`
}
