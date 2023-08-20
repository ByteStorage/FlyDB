package wal

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"os"
	"time"

	"github.com/ByteStorage/FlyDB/db/fileio"
)

const (
	// Record types
	putType    = byte(1)
	deleteType = byte(2)

	// File names for WAL
	walFileName = "/db.wal"
)

// Wal is a write-ahead log.
type Wal struct {
	m          *fileio.MMapIO // MMapIOManager
	logNum     uint32         // Log number
	saveTime   int64          // Save time
	dirPath    string         // Dir path
	readOffset int64          // Read offset
}

// NewWal creates a new WAL.
func NewWal(options Options) (*Wal, error) {
	fileName := options.DirPath + walFileName
	stat, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(options.DirPath, os.ModePerm)
			if err != nil {
				return nil, err
			}
			_, err = os.Create(fileName)
			if err != nil {
				return nil, err
			}
		}
	} else {
		options.LogNum = uint32(stat.Size() / options.FileSize)
	}
	mapIO, err := fileio.NewMMapIOManager(fileName, options.FileSize)
	if err != nil {
		return nil, err
	}
	return &Wal{
		m:        mapIO,
		logNum:   options.LogNum,
		saveTime: options.SaveTime,
		dirPath:  options.DirPath,
	}, nil
}

// Put writes a record to the WAL.
// +---------+-----------+-----------+----------------+--- ... ---+
// |CRC (4B) | Size (2B) | Type (1B) | Log number (4B)| Payload   |
// +---------+-----------+-----------+----------------+--- ... ---+
// Same as above, with the addition of
// Log number = 32bit log file number, so that we can distinguish between
// records written by the most recent log writer vs a previous one.
func (w *Wal) writeRecord(recordType byte, key, value []byte) error {
	// Prepare the payload based on record type
	var payload []byte
	switch recordType {
	case putType:
		payload = append(key, value...)
	case deleteType:
		payload = key
	default:
		return errors.New("unknown record type")
	}

	size := uint16(4 + len(payload)) // 4 bytes for log number
	buffer := make([]byte, 4+2+1+4+len(payload))

	// Compute CRC
	crc := crc32.ChecksumIEEE(buffer[4:])
	binary.LittleEndian.PutUint32(buffer, crc)

	// Write size
	binary.LittleEndian.PutUint16(buffer[4:], size)

	// Write type
	buffer[4+2] = recordType

	// Write log number
	binary.LittleEndian.PutUint32(buffer[4+2+1:], w.logNum)

	// Write payload
	copy(buffer[4+2+1+4:], payload)

	_, err := w.m.Write(buffer)
	return err
}

// Put writes a record to the WAL.
func (w *Wal) Put(key []byte, value []byte) error {
	return w.writeRecord(putType, key, value)
}

// Delete writes a delete record to the WAL.
func (w *Wal) Delete(key []byte) error {
	return w.writeRecord(deleteType, key, nil)
}

// Record is a structure that holds information about a record from the WAL.
type Record struct {
	Type  byte
	Key   []byte
	Value []byte
}

// InitReading Initializes the WAL reading position to the start of the file.
func (w *Wal) InitReading() {
	w.readOffset = 0
}

// ReadNext reads the next operation from the WAL.
func (w *Wal) ReadNext() (*Record, error) {
	buffer := make([]byte, 4+2+1+4) // Buffer size to read headers
	_, err := w.m.Read(buffer, w.readOffset)
	if err == io.EOF {
		return nil, io.EOF
	}
	if err != nil {
		return nil, err
	}

	// Move readOffset
	w.readOffset += int64(len(buffer))

	// Verify CRC
	expectedCRC := binary.LittleEndian.Uint32(buffer)
	if crc32.ChecksumIEEE(buffer[4:]) != expectedCRC {
		return nil, errors.New("corrupted record found")
	}

	// Get record size and type
	size := binary.LittleEndian.Uint16(buffer[4:])
	recordType := buffer[4+2]

	// Read the payload
	payload := make([]byte, size-4) // Subtract 4 for log number
	_, err = w.m.Read(payload, w.readOffset)
	if err != nil {
		return nil, err
	}

	// Move readOffset again
	w.readOffset += int64(len(payload))

	// Parse based on record type
	switch recordType {
	case putType:
		return &Record{Type: putType, Key: payload[:len(payload)-len(buffer)], Value: payload[len(payload)-len(buffer):]}, nil
	case deleteType:
		return &Record{Type: deleteType, Key: payload, Value: nil}, nil
	default:
		return nil, errors.New("unknown record type")
	}
}

// Save flushes the WAL to disk.
func (w *Wal) Save() error {
	return w.m.Sync()
}

// Close closes the WAL.
func (w *Wal) Close() error {
	return w.m.Close()
}

func (w *Wal) Clean() error {
	err := w.m.Close()
	if err != nil {
		return err
	}
	return os.RemoveAll(w.dirPath)
}

// AsyncSave periodically flushes the WAL to disk.
func (w *Wal) AsyncSave() {
	for range time.Tick(time.Duration(w.saveTime)) {
		err := w.Save()
		if err != nil {
			// TODO how to fix this error?
			continue
		}
	}
}
