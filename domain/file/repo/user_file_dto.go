package repo

type ClassifiedDocsQuery struct {
	UserID string
	Type   int32 // 1: pic, 2: video, 3: music, 4: document
	Offset int32
	Limit  int32
}
