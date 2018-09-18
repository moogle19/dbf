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
	dt := new(Table)

	data, err := ioutil.ReadAll(ir)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	dt.fileEncoding = fileEncoding
	dt.encoder = mahonia.NewEncoder(fileEncoding)
	dt.decoder = mahonia.NewDecoder(fileEncoding)

	header, err := parseHeader(bytes.NewReader(data[:12]))
	if err != nil {
		return nil, fmt.Errorf("failed to parse header: %v", err)
	}
	dt.Header = header

	// create fieldMap to taranslate field name to index
	dt.fieldMap = make(map[string]int)

	// Number of fields in dbase table
	dt.numberOfFields = int((dt.Header.headerSize - 1 - 32) / 32)

	// populate dbf fields
	for i := 0; i < int(dt.numberOfFields); i++ {
		offset := (i * 32) + 32

		fieldName := strings.Trim(dt.encoder.ConvertString(string(data[offset:offset+10])), "\x00")
		dt.fieldMap[fieldName] = i

		var err error

		switch data[offset+11] {
		case 'B':
			// TODO: Handle binary
			break
		case 'C':
			err = dt.AddTextField(fieldName, data[offset+16])
		case 'D':
			err = dt.AddDateField(fieldName)
		case 'N':
			err = dt.AddNumberField(fieldName, data[offset+16], data[offset+17])
		case 'L':
			err = dt.AddBooleanField(fieldName)
		case 'M':
			// TODO: Handle memo
			break
		case '@':
			// TODO: Handle Timestamp
			break
		case 'I':
			// TODO: Handle Long
			break
		case '+':
			// TODO: Handle auto increment
			break
		case 'F':
			err = dt.AddFloatField(fieldName, data[offset+16], data[offset+17])
		case 'O':
			// TODO: Handle double
			break
		case 'G':
			// TODO: Handle OLE
			break
		}

		// Check return value for errors
		if err != nil {
			return nil, err
		}
	}

	// Since we are reading dbase file from the disk at least at this
	// phase changing schema of dbase file is not allowed.
	dt.dataEntryStarted = true

	// set DbfTable dataStore slice that will store the complete file in memory
	dt.dataStore = data

	return dt, nil
}

// SaveFile saves table to a file
func (dt *Table) SaveFile(filename string) error {
	return ioutil.WriteFile(filename, append(dt.dataStore, 0x1A), os.ModePerm)
}
