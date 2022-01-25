package repo

type CreateShareDTO struct {
	UserID     string
	DocID      string
	ExpireHour int32
}
