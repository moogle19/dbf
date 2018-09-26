package dbf

import (
	"fmt"
	"strconv"
	"strings"
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

// Name returns the column name for the field
func (f *Field) Name() string {
	return f.column.Name
}

// Bool return the fields value as bool
// Default value is false
// If it is the wrong type an error is returned
func (f *Field) Bool() (bool, error) {
	if f.IsEmpty() {
		return false, nil
	}
	if f.column.Type != TypeBool {
		return false, fmt.Errorf(
			"Bool(): invalid field type: %v",
			f.column.Type,
		)
	}
	switch strings.ToLower(f.value) {
	case "t":
		return true, nil
	case "f":
		return false, nil
	default:
		return false, fmt.Errorf("Bool(): invalid value: %v", f.value)
	}
}

// Float returns the fields value as float
// If it is the wrong type or empty an error is returned
func (f *Field) Float() (float64, error) {
	if f.IsEmpty() {
		return 0.0, ErrEmptyField
	}
	if f.column.Type != TypeNumber && f.column.Type != TypeFloat {
		return 0.0, fmt.Errorf(
			"Float(): invalid field type: %v",
			f.column.Type,
		)
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
		return 0, fmt.Errorf(
			"Int(): invalid field type: %v",
			f.column.Type,
		)
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
		return 0, fmt.Errorf(
			"Int64(): invalid field type: %v",
			f.column.Type,
		)
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
		return time.Time{}, fmt.Errorf(
			"Date(): invalid field type: %v",
			f.column.Type,
		)
	}
	return time.Parse("20060102", f.value)
}
