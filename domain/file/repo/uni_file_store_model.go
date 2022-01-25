package repo

type UniFileMetaPO struct {
	Md5      string `bson:"md5,omitempty"`
	Size     int64  `bson:"size,omitempty"`
	Type     string `bson:"type,omitempty"`
	UploadBy string `bson:"upload_by,omitempty"`
}
