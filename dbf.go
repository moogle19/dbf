package dbf

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

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
	columns, err := t.parseColumns(data[32:columnHeaderSize])
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

func (c Columns) RowLength() int {
	var length int
	for _, column := range c {
		length += column.Length
	}

	return length
}

type Row struct {
	Fields []Field
}

type Field struct {
	column *Column
	value  string
}

func (f *Field) String() string {
	return f.value
}

func (f *Field) Float() (float64, error) {
	if f.column.Type != TypeNumber && f.column.Type != TypeFloat {
		return 0.0, fmt.Errorf("invalid field type")
	}
	return strconv.ParseFloat(f.value, 64)
}

func (f *Field) Int() (int, error) {
	if f.column.Type != TypeNumber {
		return 0, fmt.Errorf("invalid field type")
	}
	return strconv.Atoi(f.value)
}

func (f *Field) Date() (time.Time, error) {
	if f.column.Type != TypeDate {
		return time.Time{}, fmt.Errorf("invalid field type")
	}
	return time.Parse("20060102", f.value)
}

type Columns []*Column

type Column struct {
	Name          string
	Type          ColumnType
	Length        int
	DecimalPlaces int
}

type ColumnType rune

var (
	TypeText    ColumnType = 'C'
	TypeBool    ColumnType = 'L'
	TypeDate    ColumnType = 'D'
	TypeNumber  ColumnType = 'N'
	TypeFloat   ColumnType = 'F'
	TypeMemo    ColumnType = 'M'
	TypeUnknown ColumnType = '?'
)

func getColumnType(r byte) (ColumnType, error) {
	for _, t := range []ColumnType{TypeText, TypeBool, TypeDate, TypeNumber, TypeFloat, TypeMemo} {
		if t == ColumnType(r) {
			return ColumnType(r), nil
		}
	}
	return TypeUnknown, fmt.Errorf("column / field type %c is not supported", r)
}

func (dt *Table) parseColumns(d []byte) (Columns, error) {
	var columns []*Column
	for i := 0; i < len(d); i += 32 {
		name := strings.Trim(dt.encoder.ConvertString(string(d[i:i+10])), "\x00")
		ct, err := getColumnType(d[i+11])
		if err != nil {
			return nil, err
		}

		length := int(d[i+16])
		decimalPlaces := int(d[i+17])

		column := Column{
			Name:          name,
			Type:          ct,
			Length:        length,
			DecimalPlaces: decimalPlaces,
		}

		columns = append(columns, &column)
	}

	return columns, nil

}
