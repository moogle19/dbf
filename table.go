package dbf

import (
	"errors"
	"strconv"
	"strings"
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

// Field is a single field of the dBase table
type Field struct {
	name          string
	fieldType     string
	length        uint8
	decimalPlaces uint8
	store         [32]byte
}

// Table is a table of the dBase database
type Table struct {
	// dbase file header information
	Header *Header

	reservedBytes   [20]byte // Reserved bytes
	fieldDescriptor [32]byte // Field descriptor array
	fieldTerminator int8     // 0Dh stored as the field terminator.

	numberOfFields int // number of fiels/colums in dbase file

	// columns of dbase file
	fields []Field

	// used to map field names to index
	fieldMap map[string]int

	/*
	   "dataEntryStarted" flag is used to control whether we can change
	   dbase table structure when data enty started you can not change
	   the schema of the file if you are reading from an existing file this
	   file will be set to "true". This means you can not modify the schema
	   of a dbase table that you loaded from a file.
	*/
	dataEntryStarted bool

	// cratedFromScratch is used before adding new fields to increment nu
	createdFromScratch bool

	// encoding of dbase file
	fileEncoding string
	decoder      mahonia.Decoder
	encoder      mahonia.Encoder

	// keeps the dbase table in memory as byte array
	dataStore []byte
}

// Sets field value by index
func (dt *Table) SetFieldValueByName(row int, fieldName string, value string) error {

	i, ok := dt.fieldMap[fieldName]

	if !ok {
		return errors.New("Field name \"" + fieldName + "\" does not exist")
	}

	// set field value and return
	return dt.SetFieldValue(row, i, value)
}

// Sets field value by name
func (dt *Table) SetFieldValue(row int, fieldIndex int, value string) (err error) {

	b := []byte(dt.encoder.ConvertString(value))

	fieldLength := int(dt.fields[fieldIndex].length)

	// locate the offset of the field in DbfTable dataStore
	offset := dt.Header.HeaderSize()
	lengthOfRecord := dt.Header.RecordSize()

	offset = offset + (row * lengthOfRecord)

	recordOffset := 1

	for _, field := range dt.fields[:fieldIndex] {
		recordOffset += int(field.length)
	}

	// first fill the field with space values
	for i := 0; i < fieldLength; i++ {
		dt.dataStore[offset+recordOffset+i] = 0x20
	}

	// write new value
	switch dt.fields[fieldIndex].fieldType {
	case "C", "L", "D":
		for i := 0; i < len(b) && i < fieldLength; i++ {
			dt.dataStore[offset+recordOffset+i] = b[i]
		}
	case "N", "F":
		for i := 0; i < fieldLength; i++ {
			if i < len(b) {
				dt.dataStore[offset+recordOffset+(fieldLength-i-1)] = b[(len(b)-1)-i]
			} else {
				break
			}
		}
	}

	return
}

// RowIsDeleted checks if row is marked as deleted
func (dt *Table) RowIsDeleted(row int) bool {
	offset := int(dt.Header.headerSize)
	lengthOfRecord := int(dt.Header.recordSize)
	offset = offset + (row * lengthOfRecord)
	return (dt.dataStore[offset:(offset + 1)][0] == 0x2A)
}

// FieldValue gets the value of the specified field
func (dt *Table) FieldValue(row int, fieldIndex int) (value string) {

	offset := int(dt.Header.headerSize)
	lengthOfRecord := int(dt.Header.recordSize)

	offset = offset + (row * lengthOfRecord)

	recordOffset := 1

	for _, field := range dt.fields[:fieldIndex] {
		recordOffset += int(field.length)
	}

	temp := dt.dataStore[(offset + recordOffset):((offset + recordOffset) + int(dt.fields[fieldIndex].length))]

	for i := 0; i < len(temp); i++ {
		if temp[i] == 0x00 {
			temp = temp[0:i]
			break
		}
	}

	s := dt.decoder.ConvertString(string(temp))

	value = strings.TrimSpace(s)

	return
}

// Float64FieldValueByName retuns the value of a field given row number and fieldName provided as a float64
func (dt *Table) Float64FieldValueByName(row int, fieldName string) (value float64, err error) {

	fieldValueAsString, err := dt.FieldValueByName(row, fieldName)

	return strconv.ParseFloat(fieldValueAsString, 64)
}

// Int64FieldValueByName retuns the value of a field given row number and fieldName provided as an int64
func (dt *Table) Int64FieldValueByName(row int, fieldName string) (value int64, err error) {

	fieldValueAsString, err := dt.FieldValueByName(row, fieldName)

	return strconv.ParseInt(fieldValueAsString, 0, 64)
}

// FieldValueByName retuns the value of a field given row number and fieldName provided
func (dt *Table) FieldValueByName(row int, fieldName string) (value string, err error) {

	fieldIndex, ok := dt.fieldMap[fieldName]

	if !ok {
		err = errors.New("Field name \"" + fieldName + "\" does not exist")
		return
	}
	return dt.FieldValue(row, fieldIndex), err
}

func (dt *Table) AddNewRecord() (newRecordNumber int) {

	if dt.dataEntryStarted == false {
		dt.dataEntryStarted = true
	}

	newRecord := make([]byte, dt.Header.recordSize)
	dt.dataStore = append(dt.dataStore, newRecord...)

	// since row numbers are "0" based first we set newRecordNumber
	// and then increment number of records in dbase table
	newRecordNumber = int(dt.Header.recordCount)

	dt.Header.recordCount++
	s := uint32ToBytes(dt.Header.recordCount)
	dt.dataStore[4] = s[0]
	dt.dataStore[5] = s[1]
	dt.dataStore[6] = s[2]
	dt.dataStore[7] = s[3]

	return newRecordNumber
}

func (dt *Table) AddTextField(fieldName string, length uint8) (err error) {
	return dt.addField(fieldName, 'C', length, 0)
}

func (dt *Table) AddBooleanField(fieldName string) (err error) {
	return dt.addField(fieldName, 'L', 1, 0)
}

func (dt *Table) AddDateField(fieldName string) (err error) {
	return dt.addField(fieldName, 'D', 8, 0)
}

func (dt *Table) AddNumberField(fieldName string, length uint8, decimalPlaces uint8) (err error) {
	return dt.addField(fieldName, 'N', length, decimalPlaces)
}

func (dt *Table) AddFloatField(fieldName string, length uint8, decimalPlaces uint8) (err error) {
	return dt.addField(fieldName, 'F', length, decimalPlaces)
}

// NumberOfRecords return number of rows in dbase table
func (dt *Table) NumberOfRecords() int {
	return int(dt.Header.recordCount)
}

// Fields return slice of DbfField
func (dt *Table) Fields() []Field {
	return dt.fields
}

// FieldNames return slice of DbfField names
func (dt *Table) FieldNames() []string {
	names := make([]string, 0)

	for _, field := range dt.Fields() {
		names = append(names, field.name)
	}

	return names
}

func (dt *Table) addField(fieldName string, fieldType byte, length uint8, decimalPlaces uint8) (err error) {

	if dt.dataEntryStarted {
		return errors.New("Once you start entering data to the dbase table or open an existing dbase file, altering dbase table schema is not allowed!")
	}

	normalizedFieldName := dt.getNormalizedFieldName(fieldName)

	if dt.HasField(normalizedFieldName) {
		return errors.New("Field name \"" + normalizedFieldName + "\" already exists!")
	}

	df := new(Field)
	df.name = normalizedFieldName
	df.fieldType = string(fieldType)
	df.length = length
	df.decimalPlaces = decimalPlaces

	slice := dt.convertToByteSlice(normalizedFieldName, 10)

	//fmt.Printf("len slice:%v\n", len(slice))

	// Field name in ASCII (max 10 chracters)
	for i := 0; i < len(slice); i++ {
		df.store[i] = slice[i]
		//fmt.Printf("i:%s\n", string(slice[i]))
	}

	// Field names are terminated by 00h
	df.store[10] = 0x00

	// Set field's data type
	// C (Character)  All OEM code page characters.
	// D (Date)     Numbers and a character to separate month, day, and year (stored internally as 8 digits in YYYYMMDD format).
	// N (Numeric)    - . 0 1 2 3 4 5 6 7 8 9
	// F (Floating Point)   - . 0 1 2 3 4 5 6 7 8 9
	// L (Logical)    ? Y y N n T t F f (? when not initialized).
	df.store[11] = fieldType

	// length of field
	df.store[16] = length

	// number of decimal places
	// Applicable only to number/float
	df.store[17] = df.decimalPlaces

	dt.fields = append(dt.fields, *df)

	// if createdFromScratch we need to update dbase header to reflect the changes we have made
	if dt.createdFromScratch {
		dt.updateHeader()
	}

	return
}

// updateHeader updates the dbase file header after a field added
func (dt *Table) updateHeader() {
	// first create a slice from initial 32 bytes of datastore as the foundation of the new slice
	// later we will set this slice to dt.dataStore to create the new header slice
	slice := dt.dataStore[0:32]

	// set dbase file signature
	slice[0] = 0x03

	var lengthOfEachRecord uint16

	for i, field := range dt.Fields() {
		lengthOfEachRecord += uint16(field.length)
		slice = append(slice, field.store[:]...)

		// don't forget to update fieldMap. We need it to find the index of a field name
		dt.fieldMap[field.name] = i
	}

	// end of file header terminator (0Dh)
	slice = append(slice, 0x0D)

	// now reset dt.dataStore slice with the updated one
	dt.dataStore = slice

	// update the number of bytes in dbase file header
	dt.Header.headerSize = uint16(len(slice))
	s := uint32ToBytes(uint32(dt.Header.headerSize))
	dt.dataStore[8] = s[0]
	dt.dataStore[9] = s[1]

	dt.Header.recordSize = lengthOfEachRecord + 1 // dont forget to add "1" for deletion marker which is 20h

	// update the lenght of each record
	s = uint32ToBytes(uint32(dt.Header.recordSize))
	dt.dataStore[10] = s[0]
	dt.dataStore[11] = s[1]

	return
}

func (dt *Table) GetRowAsSlice(row int) []string {

	s := make([]string, len(dt.Fields()))

	for i := 0; i < len(dt.Fields()); i++ {
		s[i] = dt.FieldValue(row, i)
	}

	return s
}

func (dt *Table) HasField(fieldName string) bool {

	for i := 0; i < len(dt.fields); i++ {
		if dt.fields[i].name == fieldName {
			return true
		}
	}

	return false
}

func (dt *Table) DecimalPlacesInField(fieldName string) (uint8, error) {
	if !dt.HasField(fieldName) {
		return 0, errors.New("Field name \"" + fieldName + "\" does not exist. ")
	}

	for i := 0; i < len(dt.fields); i++ {
		if dt.fields[i].name == fieldName {
			if dt.fields[i].fieldType == "N" || dt.fields[i].fieldType == "F" {
				return dt.fields[i].decimalPlaces, nil
			}
		}
	}

	return 0, errors.New("Type of field \"" + fieldName + "\" is not Numeric or Float.")
}

// convertToBytesSlice converts value to byte slice according to given encoding and return
// a slice that is length equals to numberOfBytes or less if the string is shorter than
// numberOfBytes
func (dt *Table) convertToByteSlice(value string, numberOfBytes int) (s []byte) {
	e := mahonia.NewEncoder(dt.fileEncoding)
	b := []byte(e.ConvertString(value))

	if len(b) <= numberOfBytes {
		s = b
	} else {
		s = b[0:numberOfBytes]
	}
	return
}

func (dt *Table) getNormalizedFieldName(name string) (s string) {
	e := mahonia.NewEncoder(dt.fileEncoding)
	b := []byte(e.ConvertString(name))

	if len(b) > 10 {
		b = b[0:10]
	}

	d := mahonia.NewDecoder(dt.fileEncoding)
	s = d.ConvertString(string(b))

	return
}

// New creates a new dBase database from scratch
func New(encoding string) (table *Table) {

	// Create and populate DbaseTable struct
	dt := new(Table)

	dt.fileEncoding = encoding
	dt.encoder = mahonia.NewEncoder(encoding)
	dt.decoder = mahonia.NewDecoder(encoding)

	// set whether or not this table has been created from scratch
	dt.createdFromScratch = true

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

	// create fieldMap to taranslate field name to index
	dt.fieldMap = make(map[string]int)

	// Number of fields in dbase table
	dt.Header.recordCount = uint32((dt.Header.headerSize - 1 - 32) / 32)

	s := make([]byte, dt.Header.headerSize)

	// Since we are reading dbase file from the disk at least at this
	// phase changing schema of dbase file is not allowed.
	dt.dataEntryStarted = false

	// set DbfTable dataStore slice that will store the complete file in memory
	dt.dataStore = s

	dt.dataStore[0] = dt.Header.Signature
	dt.dataStore[1] = dt.Header.updateYear
	dt.dataStore[2] = dt.Header.updateMonth
	dt.dataStore[3] = dt.Header.updateDay

	// no MDX file (index upon demand)
	dt.dataStore[28] = 0x00

	// set dbase language driver
	// Huston we have problem!
	// There is no easy way to deal with encoding issues. At least at the moment
	// I will try to find archaic encoding code defined by dbase standard (if there is any)
	// for given encoding. If none match I will go with default ANSI.
	//
	// Despite this flag in set in dbase file, I will continue to use provide encoding for
	// the everything except this file encoding flag.
	//
	// Why? To make sure at least if you know the real encoding you can process text accordingly.

	if code, ok := encodingTable[lookup[encoding]]; ok {
		dt.dataStore[28] = code
	} else {
		dt.dataStore[28] = 0x57 // ANSI
	}

	return dt
}