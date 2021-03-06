package dbf

import "fmt"

// ColumnType is the data type of the dbf column
type ColumnType rune

func (c ColumnType) String() string {
	switch c {
	case TypeText:
		return "Text"
	case TypeBool:
		return "Bool"
	case TypeDate:
		return "Date"
	case TypeNumber:
		return "Number"
	case TypeFloat:
		return "Float"
	case TypeMemo:
		return "Memo"
	default:
		return "unknown: " + string(c)
	}
}

var (
	// TypeText is a string field
	TypeText ColumnType = 'C'
	// TypeBool is a boolean field
	TypeBool ColumnType = 'L'
	// TypeDate is a date field
	TypeDate ColumnType = 'D'
	// TypeNumber is an integer number
	TypeNumber ColumnType = 'N'
	// TypeFloat is a float number
	TypeFloat ColumnType = 'F'
	// TypeMemo is a memo
	TypeMemo ColumnType = 'M'
	// TypeUnknown is used when the type is not known
	TypeUnknown ColumnType = '?'
)

func getColumnType(r byte) (ColumnType, error) {

	allowedTypes := []ColumnType{
		TypeText,
		TypeBool,
		TypeDate,
		TypeNumber,
		TypeFloat,
		TypeMemo,
	}

	for _, t := range allowedTypes {
		if t == ColumnType(r) {
			return ColumnType(r), nil
		}
	}
	return TypeUnknown, fmt.Errorf("column / field type %c is not supported", r)
}
