package repo

import (
	"encoding/json"
	"time"

	"github.com/changpro/disk-service/infra/constants"
)

type RecycleFilePO struct {
	ID       string    `bson:"_id,omitempty" json:"docId,omitempty"`
	UserID   string    `bson:"user_id,omitempty" json:"userId,omitempty"`
	Name     string    `bson:"name,omitempty" json:"docName,omitempty"`
	IsDir    int32     `bson:"is_dir,omitempty" json:"isDir,omitempty"` // 2: file, 1: directory
	DeleteAt time.Time `bson:"delete_at,omitempty" json:"-"`
	ExpireAt time.Time `bson:"expire_at,omitempty" json:"-"`
}

func (p *RecycleFilePO) MarshalJSON() ([]byte, error) {
	type Alias RecycleFilePO
	return json.Marshal(&struct {
		*Alias
		DeleteAt string `json:"deleteAt"`
		ExpireAt string `json:"expireAt"`
	}{
		Alias:    (*Alias)(p),
		DeleteAt: p.DeleteAt.Format(constants.StandardTimeFormat),
		ExpireAt: p.ExpireAt.Format(constants.StandardTimeFormat),
	})
}
