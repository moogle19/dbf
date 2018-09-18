package dbf

import (
	"fmt"
	"strconv"
	"time"
)

// Field represents a single field in the the dbf table
type Field struct {
	column *Column
	value  string
}

// IsEmpty checks if a field is empty
func (f *Field) IsEmpty() bool {
	return f == nil || f.value == ""
}

// String returns the fields value as string
func (f *Field) String() string {
	if f == nil {
		return ""
	}
	return f.value
}

// Float returns the fields value as float
// If it is the wrong type or empty an error is returned
func (f *Field) Float() (float64, error) {
	if f.IsEmpty() {
		return 0.0, ErrEmptyField
	}
	if f.column.Type != TypeNumber && f.column.Type != TypeFloat {
		return 0.0, fmt.Errorf("invalid field type")
	}
	return strconv.ParseFloat(f.value, 64)
}

// Int returns the fields value as int
// If it is the wrong type or empty an error is returned
func (f *Field) Int() (int, error) {
	if f.IsEmpty() {
		return 0, ErrEmptyField
	}
	if f.column.Type != TypeNumber {
		return 0, fmt.Errorf("invalid field type")
	}

	return strconv.Atoi(f.value)
}

// Int64 returns the fields value as int64
// If it is the wrong type or empty an error is returned
func (f *Field) Int64() (int64, error) {
	if f.IsEmpty() {
		return 0, ErrEmptyField
	}
	if f.column.Type != TypeNumber {
		return 0, fmt.Errorf("invalid field type")
	}

	return strconv.ParseInt(f.value, 10, 64)
}

// Date returns the fields value as time.Time
// If it is the wrong type or empty an error is returned
func (f *Field) Date() (time.Time, error) {
	if f.IsEmpty() {
		return time.Time{}, ErrEmptyField
	}
	if f.column.Type != TypeDate {
		return time.Time{}, fmt.Errorf("invalid field type")
	}
	return time.Parse("20060102", f.value)
}
