package memory

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
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

type Wal struct {
	m        *fileio.MMapIO
	logNum   uint32
	saveTime int64
}

func NewWal(options Options) (*Wal, error) {
	mapIO, err := fileio.NewMMapIOManager(options.Option.DirPath+walFileName, options.FileSize)
	if err != nil {
		return nil, err
	}
	return &Wal{
		m:        mapIO,
		logNum:   options.LogNum,
		saveTime: options.SaveTime,
	}, nil
}

// Put writes a record to the WAL.
//	+---------+-----------+-----------+----------------+--- ... ---+
//	|CRC (4B) | Size (2B) | Type (1B) | Log number (4B)| Payload   |
//	+---------+-----------+-----------+----------------+--- ... ---+
//	Same as above, with the addition of
//	Log number = 32bit log file number, so that we can distinguish between
//	records written by the most recent log writer vs a previous one.
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

func (w *Wal) Save() error {
	return w.m.Sync()
}

func (w *Wal) Close() error {
	return w.m.Close()
}

func (w *Wal) AsyncSave() {
	for range time.Tick(time.Duration(w.saveTime)) {
		err := w.Save()
		if err != nil {
			// TODO how to fix this error?
			continue
		}
	}
}
