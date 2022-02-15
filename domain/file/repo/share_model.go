package repo

type ShareDetailPO struct {
	Uploader    string `redis:"uploader" json:"uploader,omitempty"`
	Password    string `redis:"password" json:"password,omitempty"`
	DocID       string `redis:"docId" json:"docId,omitempty"`
	DocName     string `redis:"docName" json:"docName,omitempty"`
	UniFileID   string `redis:"uniFileId" json:"uniFileId"`
	DocSize     int64  `redis:"docSize" json:"docSize,omitempty"`
	DocType     string `redis:"docType" json:"docType,omitempty"`
	IsDir       int32  `redis:"isDir" json:"isDir,omitempty"`     // 1 folder, 2 file
	FileNum     int32  `redis:"fileNum" json:"fileNum,omitempty"` // specificly, if it's a file, this field has a value of 1;
	CreateTime  string `redis:"createTime" json:"createTime,omitempty"`
	ExpireHours int32  `redis:"expireHours" json:"expireHours,omitempty"`
	ViewNum     int32  `redis:"viewNum" json:"viewNum,omitempty"`
	SaveNum     int32  `redis:"saveNum" json:"saveNum,omitempty"`
}
