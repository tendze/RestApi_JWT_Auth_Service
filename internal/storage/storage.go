package storage

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidLoginOrPass = errors.New("invalid login or password")
	ErrUserExists         = errors.New("user already exists")
)
