package repo

type ShareDetailPO struct {
	Uploader    string `redis:"uploader" json:"uploader,omitempty"`
	DocID       string `redis:"docId" json:"docId,omitempty"`
	DocName     string `redis:"docName" json:"docName,omitempty"`
	DocSize     int64  `redis:"docSize" json:"docSize,omitempty"`
	DocType     int32  `redis:"docType" json:"docType,omitempty"` // 1 folder, 2 file
	FileNum     int32  `redis:"fileNum" json:"fileNum,omitempty"` // specificly, if it's a file, this field has a value of 1;
	CreateTime  string `redis:"createTime" json:"createTime,omitempty"`
	ExpireHours int32  `redis:"expireHours" json:"expireHours,omitempty"`
	ViewNum     int32  `redis:"viewNum" json:"viewNum,omitempty"`
	SaveNum     int32  `redis:"saveNum" json:"saveNum,omitempty"`
}