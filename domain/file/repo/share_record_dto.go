package repo

type RecordQuery struct {
	UserID    string
	Offset    int32
	Limit     int32
	Type      int32
	StartTime int64
	EndTime   int64
}
