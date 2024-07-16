package services

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrNoFreeSlot   = errors.New("no free slot")
)
