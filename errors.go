package dbf

import "errors"

var (
	// ErrInvalidColumnData is returned when the column data is invalid
	ErrInvalidColumnData = errors.New("invalid column data")
	// ErrEmptyField is returned when a field is empty
	ErrEmptyField = errors.New("field is empty")
	// ErrInvalidFieldName is returned when a given field name was not found
	ErrInvalidFieldName = errors.New("invalid field name")
	// ErrIndexOutOfBounds is returned when a row index is out of bounds
	ErrIndexOutOfBounds = errors.New("index is out of bounds")
)
