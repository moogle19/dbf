package dbf

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/text/encoding"
)

// Open opens a DBF Table from an io.Reader
func Open(r io.Reader) (*Table, error) {
	return OpenWithEncoding(r, encoding.Nop)
}

func OpenWithEncoding(r io.Reader, enc encoding.Encoding) (*Table, error) {
	return createDbfTable(r, enc)
}

// OpenFile opens a DBF Table from file
func OpenFile(filename string) (*Table, error) {
	return OpenFileWithEncoding(filename, encoding.Nop)
}

func OpenFileWithEncoding(filename string, enc encoding.Encoding) (*Table, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return OpenWithEncoding(f, enc)
}

func createDbfTable(ir io.Reader, enc encoding.Encoding) (table *Table, err error) {
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

	columns, err := parseColumns(data[32:columnHeaderSize], 32, enc)
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
		row, err := parseRow(rowData[i:i+rowLength+1], columns, enc)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)

	}

	t.Columns = columns
	t.Rows = rows

	return t, nil
}
