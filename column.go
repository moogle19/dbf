package dbf

import (
	"bytes"

	"golang.org/x/text/encoding/charmap"
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

func newColumn(rawData []byte, encoding Encoding) (*Column, error) {
	if len(rawData) != 32 {
		return nil, ErrInvalidColumnData
	}

	nameData := rawData[:10]
	if encoding != nil {
		var err error
		dec := (*charmap.Charmap)(encoding).NewDecoder()
		nameData, err = dec.Bytes(rawData[:10])
		if err != nil {
			return nil, err
		}
	}
	name := string(bytes.Trim(nameData, "\x00"))

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

func parseColumns(rawData []byte, columnLength int, encoding Encoding) (Columns, error) {
	var columns []*Column

	for i := 0; i < len(rawData); i += columnLength {

		column, err := newColumn(rawData[i:i+columnLength], encoding)
		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	return columns, nil

}
