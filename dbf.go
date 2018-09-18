package dbf

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/axgle/mahonia"
)

// Open opens a DBF Table from an io.Reader
func Open(r io.Reader, encoding string) (*Table, error) {
	return createDbfTable(r, encoding)
}

// OpenFile opens a DBF Table from file
func OpenFile(filename string, encoding string) (*Table, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Open(f, encoding)
}

func createDbfTable(ir io.Reader, fileEncoding string) (table *Table, err error) {
	// Create and pupulate DbaseTable struct
	t := new(Table)

	// read complete table
	data, err := ioutil.ReadAll(ir)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Initalize encoder
	t.fileEncoding = fileEncoding
	t.encoder = mahonia.NewEncoder(fileEncoding)
	t.decoder = mahonia.NewDecoder(fileEncoding)

	// Parse header
	header, err := parseHeader(bytes.NewReader(data[:12]))
	if err != nil {
		return nil, fmt.Errorf("failed to parse header: %v", err)
	}
	t.Header = header

	// Parse columns
	fieldCount := (t.Header.HeaderSize() - 1) / 32

	columnHeaderSize := fieldCount * 32

	columns, err := parseColumns(data[32:columnHeaderSize], 32, &t.decoder)
	if err != nil {
		return nil, err
	}

	// Parse rows
	offset := int(header.HeaderSize())

	rowData := data[offset+1:]
	rowLength := columns.RowLength()

	var rows []*Row

	// Iterate rows
	for i := 0; i < len(data); i += rowLength + 1 {
		var r Row

		off := i
		for _, c := range columns {
			var field Field
			field.column = c

			l := c.Length
			for i, b := range rowData[off : off+c.Length] {
				if b == byte(0) {
					l = i
					break
				}
			}
			field.value = strings.TrimSpace(t.decoder.ConvertString(string(rowData[off : off+l])))
			r.Fields = append(r.Fields, field)

			off += c.Length
		}

		rows = append(rows, &r)

	}

	t.Columns = columns
	t.Rows = rows

	return t, nil
}
