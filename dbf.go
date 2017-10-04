package dbf

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/axgle/mahonia"
)

func NewFromFile(fileName string, fileEncoding string) (table *Table, err error) {
	s, err := readFile(fileName)
	if err != nil {
		return nil, err
	}

	return createDbfTable(s, fileEncoding)
}

func NewFromByteArray(data []byte, fileEncoding string) (table *Table, err error) {
	return createDbfTable(data, fileEncoding)
}

type meta struct {
	signature   uint8
	updateYear  uint8
	updateMonth uint8
	updateDay   uint8
	recordCount uint32
	headerSize  uint16
	recordSize  uint16
}

func parseMetadata(reader io.Reader) (*meta, error) {
	m := make([]byte, 11)
	n, err := reader.Read(m)
	if err != nil {
		return nil, err
	} else if n != 11 {
		return nil, fmt.Errorf("file too short: %d bytes", n)
	}

	return &meta{
		signature:   m[0],
		updateYear:  m[1],
		updateMonth: m[2],
		updateDay:   m[3],
		recordCount: uint32(m[4]) | (uint32(m[5]) << 8) | (uint32(m[6]) << 16) | (uint32(m[7]) << 24),
		headerSize:  uint16(m[8]) | (uint16(m[9]) << 8),
		recordSize:  uint16(m[10]) | (uint16(m[11]) << 8),
	}, nil
}

func createDbfTable(s []byte, fileEncoding string) (table *Table, err error) {
	// Create and pupulate DbaseTable struct
	dt := new(Table)

	dt.fileEncoding = fileEncoding
	dt.encoder = mahonia.NewEncoder(fileEncoding)
	dt.decoder = mahonia.NewDecoder(fileEncoding)

	// read dbase table header information
	dt.fileSignature = s[0]
	dt.updateYear = s[1]
	dt.updateMonth = s[2]
	dt.updateDay = s[3]
	dt.numberOfRecords = uint32(s[4]) | (uint32(s[5]) << 8) | (uint32(s[6]) << 16) | (uint32(s[7]) << 24)
	dt.numberOfBytesInHeader = uint16(s[8]) | (uint16(s[9]) << 8)
	dt.lengthOfEachRecord = uint16(s[10]) | (uint16(s[11]) << 8)

	// create fieldMap to taranslate field name to index
	dt.fieldMap = make(map[string]int)

	// Number of fields in dbase table
	dt.numberOfFields = int((dt.numberOfBytesInHeader - 1 - 32) / 32)

	// populate dbf fields
	for i := 0; i < int(dt.numberOfFields); i++ {
		offset := (i * 32) + 32

		fieldName := strings.Trim(dt.encoder.ConvertString(string(s[offset:offset+10])), string([]byte{0}))
		dt.fieldMap[fieldName] = i

		var err error

		switch s[offset+11] {
		case 'C':
			err = dt.AddTextField(fieldName, s[offset+16])
		case 'N':
			err = dt.AddNumberField(fieldName, s[offset+16], s[offset+17])
		case 'F':
			err = dt.AddFloatField(fieldName, s[offset+16], s[offset+17])
		case 'L':
			err = dt.AddBooleanField(fieldName)
		case 'D':
			err = dt.AddDateField(fieldName)
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
	dt.dataStore = s

	return dt, nil
}

func (dt *Table) SaveFile(filename string) (err error) {

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	dsBytes, dsErr := f.Write(dt.dataStore)

	if dsErr != nil {
		return dsErr
	}

	// Add dbase end of file marker (1Ah)

	footerByte, footerErr := f.Write([]byte{0x1A})

	if footerErr != nil {
		return footerErr
	}

	fmt.Printf("%v bytes written to file '%v'.\n", dsBytes+footerByte, filename)

	return nil
}
