package entity

import "errors"

var (
	ErrUserExists   = errors.New("user with this login already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidCreds = errors.New("invalid login or password")
	
	ErrAdNotFound = errors.New("ad not found")
)
