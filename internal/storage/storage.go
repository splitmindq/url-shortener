package storage

import "errors"

var (
	ErrUrlNotFound      = errors.New("URL not found")
	ErrUrlAlreadyExists = errors.New("URL already exists")
)
