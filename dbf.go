package dbf

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// Open opens a DBF Table from an io.Reader
func Open(r io.Reader) (*Table, error) {
	return createDbfTable(r)
}

// OpenFile opens a DBF Table from file
func OpenFile(filename string) (*Table, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Open(f)
}

func createDbfTable(ir io.Reader) (table *Table, err error) {
	// Create and pupulate DbaseTable struct
	t := new(Table)

	// read complete table
	data, err := ioutil.ReadAll(ir)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Parse header
	header, err := parseHeader(bytes.NewReader(data[:12]))
	if err != nil {
		return nil, fmt.Errorf("failed to parse header: %v", err)
	}
	t.Header = header

	// Parse columns
	fieldCount := (t.Header.HeaderSize() - 1) / 32

	columnHeaderSize := fieldCount * 32

	columns, err := parseColumns(data[32:columnHeaderSize], 32)
	if err != nil {
		return nil, err
	}

	// Parse rows
	offset := int(header.HeaderSize())

	rowData := data[offset+1:]
	rowLength := columns.RowLength()

	var rows []*Row

	// Iterate rows
	for i := 0; i < len(rowData); i += rowLength + 1 {
		row, err := parseRow(rowData[i:i+rowLength+1], columns)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)

	}

	t.Columns = columns
	t.Rows = rows

	return t, nil
}
