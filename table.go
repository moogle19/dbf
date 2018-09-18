package dbf

import (
	"time"

	"github.com/axgle/mahonia"
)

const (
	// Version5 represents dBase Level5
	Version5 Version = 3
	// Version7 represents dBase Level7
	Version7 Version = 4
	// VersionUnknown represents an unknown dBase Level
	VersionUnknown Version = 0
)

// Version defines the dBase Level / Version
type Version byte

func (v Version) String() string {
	switch v {
	case Version5:
		return "dBase Level 5"
	case Version7:
		return "dBase Level 7"
	default:
		return "unknown dBase Level"
	}
}

// Table is a table of the dBase database
type Table struct {
	// dbase file header information
	Header  *Header
	Columns Columns
	Rows    []*Row

	fileEncoding string
	encoder      mahonia.Encoder
	decoder      mahonia.Decoder
}

// New creates a new dBase database from scratch
func New(encoding string) (table *Table) {

	// Create and populate DbaseTable struct
	dt := new(Table)

	dt.fileEncoding = encoding
	dt.encoder = mahonia.NewEncoder(encoding)
	dt.decoder = mahonia.NewDecoder(encoding)

	// set whether or not this table has been created from scratch

	// read dbase table header information
	dt.Header = &Header{
		Signature:   0x03,
		updateYear:  byte(time.Now().Year() % 100),
		updateMonth: byte(time.Now().Month()),
		updateDay:   byte(time.Now().YearDay()),
		recordCount: 0,
		headerSize:  32,
		recordSize:  0,
	}

	return dt
}
