package dbf

import "errors"

var (
	ErrInvalidColumnData = errors.New("invalid column data")
	ErrEmptyField        = errors.New("field is empty")
	ErrInvalidFieldName  = errors.New("invalid field name")
	ErrIndexOutOfBounds  = errors.New("index is out of bounds")
)
