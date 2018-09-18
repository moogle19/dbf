package dbf

import (
	"strings"

	"github.com/axgle/mahonia"
)

// Column represents a dBase column
type Column struct {
	Name          string
	Type          ColumnType
	Length        int
	DecimalPlaces int
	index         int
}

// Columns is a slice of Columns
type Columns []*Column

func newColumn(rawData []byte, dec *mahonia.Decoder) (*Column, error) {
	if len(rawData) != 32 {
		return nil, ErrInvalidColumnData
	}

	name := strings.Trim(dec.ConvertString(string(rawData[:10])), "\x00")
	ct, err := getColumnType(rawData[11])
	if err != nil {
		return nil, err
	}

	length := int(rawData[16])
	decimalPlaces := int(rawData[17])

	return &Column{
		Name:          name,
		Type:          ct,
		Length:        length,
		DecimalPlaces: decimalPlaces,
	}, nil

}

// RowLength returns the length of a row
func (c Columns) RowLength() int {
	var length int
	for _, column := range c {
		length += column.Length
	}

	return length
}

func parseColumns(rawData []byte, columnLength int, dec *mahonia.Decoder) (Columns, error) {
	var columns []*Column

	for i := 0; i < len(rawData); i += columnLength {

		column, err := newColumn(rawData[i:i+columnLength], dec)
		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	return columns, nil

}
