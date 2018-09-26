package dbf

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type Header struct {
	Version     Version
	Signature   uint8
	updateYear  uint8
	updateMonth uint8
	updateDay   uint8
	recordCount uint32
	headerSize  uint16
	recordSize  uint16
}

func (h *Header) RecordCount() int {
	return int(h.recordCount)
}

func (h *Header) HeaderSize() int {
	return int(h.headerSize)
}

func (h *Header) RecordSize() int {
	return int(h.recordSize)
}

func (h *Header) UpdatedAt() time.Time {
	return time.Date(int(h.updateYear), time.Month(h.updateMonth), int(h.updateDay), 0, 0, 0, 0, time.Local)
}

func parseHeader(reader io.Reader) (*Header, error) {
	m := make([]byte, 12)
	n, err := reader.Read(m)
	if err != nil {
		return nil, err
	} else if n != 12 {
		return nil, fmt.Errorf("file too short: %d bytes", n)
	}

	var version Version
	switch Version(m[0] & 0x7) {
	case Version5:
		version = Version5
	case Version7:
		version = Version7
	default:
		version = VersionUnknown
	}

	return &Header{
		Version:     version,
		Signature:   m[0],
		updateYear:  m[1],
		updateMonth: m[2],
		updateDay:   m[3],
		recordCount: binary.LittleEndian.Uint32(m[4:8]),
		headerSize:  binary.LittleEndian.Uint16(m[8:10]),
		recordSize:  binary.LittleEndian.Uint16(m[10:12]),
	}, nil
}
