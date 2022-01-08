package service

import "context"

type UserDTO struct {
	UserID       string
	PasswordMask string
}

func RegisterNewUser(ctx context.Context, user UserDTO) error {
	return nil
}

func SignIn(ctx context.Context, user UserDTO) error {
	return nil
}
