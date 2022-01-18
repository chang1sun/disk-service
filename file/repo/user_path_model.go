package repo

type UserPathPO struct {
	ID     string `bson:"_id,omitempty" json:"id,omitempty"`
	UserID string `bson:"user_id,omitempty" json:"userId,omitempty"`
	Path   string `bson:"path,omitempty" json:"path,omitempty"`
}
