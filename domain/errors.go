package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("your requested data is not found")
	ErrConflict            = errors.New("your data already exist")
	ErrBadParamInput       = errors.New("given param is not valid")
)
