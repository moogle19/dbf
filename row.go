package dbf

import (
	"fmt"
	"strings"

	"github.com/axgle/mahonia"
)

// Row represents a single row in the dbf database
type Row struct {
	fields map[string]*Field
}

func (r *Row) String() string {
	str := ""

	for k, v := range r.fields {
		if str == "" {
			str += "["
		} else {
			str += ", "
		}
		str += fmt.Sprintf("%s -> %v", k, v)
	}

	return str + "]"
}

func parseRow(rawData []byte, columns Columns, dec *mahonia.Decoder) (*Row, error) {
	r := newRow()

	var offset int

	for _, c := range columns {

		if offset > len(rawData) {
			return nil, ErrIndexOutOfBounds
		}

		length := c.Length
		for i, b := range rawData[offset : offset+c.Length] {
			if b == byte(0) {
				length = i
				break
			}
		}

		value := strings.TrimSpace(dec.ConvertString(string(rawData[offset : offset+length])))

		r.fields[c.Name] = &Field{
			column: c,
			value:  value,
		}

		offset += c.Length
	}

	return r, nil
}

func newRow() *Row {
	return &Row{
		fields: make(map[string]*Field),
	}
}

// IsEmpty checks if a row is empty
func (r *Row) IsEmpty() bool {
	return len(r.fields) == 0
}

// FieldByName returns the field with the specified name
// If not found, nil is returned
func (r *Row) FieldByName(name string) (*Field, error) {
	val, ok := r.fields[name]
	if !ok {
		return nil, ErrInvalidFieldName
	}
	return val, nil
}

func (r *Row) FieldByIndex(index int) (*Field, error) {
	if index >= len(r.fields) {
		return nil, ErrIndexOutOfBounds
	}

	for _, field := range r.fields {
		if field.column.index == index {
			return field, nil
		}
	}
	return nil, ErrIndexOutOfBounds
}
