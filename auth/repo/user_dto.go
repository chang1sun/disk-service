package repo

type ModifyPwDTO struct {
	UserID    string
	AuthEmail string
	OldPw     string
	NewPw     string
}

type ModifyUserProfileDTO struct {
	UserID    string
	AuthEmail string
	Icon      string
}
