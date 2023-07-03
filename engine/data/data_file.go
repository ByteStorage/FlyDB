package data

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/engine/fileio"
	"hash/crc32"
	"io"
	"path/filepath"
)

const (
	DataFileSuffix      = ".data"
	HintFileSuffix      = "hintIndex"
	MergeFinaFileSuffix = "mergeFina"
)

// DataFile represents a data file.
type DataFile struct {
	FileID    uint32           // File ID
	WriteOff  int64            // Position where the file is currently being written
	IoManager fileio.IOManager // IO read/write operations
}

// OpenDataFile opens a new data file.
func OpenDataFile(dirPath string, fileID uint32, fileSize int64, fioType int8) (*DataFile, error) {
	fileName := GetDataFileName(dirPath, fileID)
	return newDataFile(fileName, fileID, fileSize, fioType)
}

// GetDataFileName returns the file name for a data file.
func GetDataFileName(dirPath string, fileID uint32) string {
	return filepath.Join(dirPath, fmt.Sprintf("%09d", fileID)+DataFileSuffix)
}

// OpenHintFile opens the hint index file.
func OpenHintFile(dirPath string, fileSize int64, fioType int8) (*DataFile, error) {
	fileName := filepath.Join(dirPath, HintFileSuffix)
	return newDataFile(fileName, 0, fileSize, fioType)
}

// OpenMergeFinaFile opens the file that indicates merge completion.
func OpenMergeFinaFile(dirPath string, fileSize int64, fioType int8) (*DataFile, error) {
	fileName := filepath.Join(dirPath, MergeFinaFileSuffix)
	return newDataFile(fileName, 0, fileSize, fioType)
}

func newDataFile(dirPath string, fileID uint32, fileSize int64, fioType int8) (*DataFile, error) {
	// Initialize the IOManager interface
	ioManager, err := fileio.NewIOManager(dirPath, fileSize, fioType)
	if err != nil {
		return nil, err
	}
	return &DataFile{
		FileID:    fileID,
		WriteOff:  0,
		IoManager: ioManager,
	}, nil
}

// ReadLogRecord reads a log record from the data file based on the offset.
func (df *DataFile) ReadLogRecord(offset int64) (*LogRecord, int64, error) {
	fileSize, err := df.IoManager.Size()
	if err != nil {
		return nil, 0, err
	}

	var headerBytes int64 = maxLogRecordHeaderSize
	if offset+maxLogRecordHeaderSize > fileSize {
		headerBytes = fileSize - offset
	}

	// Read header information
	headerBuf, err := df.readNBytes(headerBytes, offset)
	if err != nil {
		return nil, 0, err
	}

	header, headerSize := decodeLogRecordHeader(headerBuf)
	// The following conditions indicate reaching the end of the file, directly return EOF
	if header == nil {
		return nil, 0, io.EOF
	}
	if header.crc == 0 && header.keySize == 0 && header.valueSize == 0 {
		return nil, 0, io.EOF
	}

	// Retrieve the lengths of the key and value
	keySize, valueSize := int64(header.keySize), int64(header.valueSize)
	var recordSize = headerSize + keySize + valueSize

	logRecord := &LogRecord{Type: header.recordType}

	// Read the actual user-stored key/value data
	if keySize > 0 || valueSize > 0 {
		kvBuf, err := df.readNBytes(keySize+valueSize, headerSize+offset)
		if err != nil {
			return nil, 0, err
		}
		// Decode
		logRecord.Key = kvBuf[:keySize]
		logRecord.Value = kvBuf[keySize:]
	}

	// Verify CRC (check data integrity)
	crc := getLogRecordCRC(logRecord, headerBuf[crc32.Size:headerSize])
	if crc != header.crc {
		return nil, 0, ErrInvalidCRC
	}
	return logRecord, recordSize, nil
}

func (df *DataFile) Write(buf []byte) error {
	size, err := df.IoManager.Write(buf)
	if err != nil {
		return err
	}
	df.WriteOff += int64(size)
	return nil
}

// WriteHintRecord writes index information to the hint file.
func (df *DataFile) WriteHintRecord(key []byte, pst *LogRecordPst) error {
	record := &LogRecord{
		Key:   key,
		Value: EncodeLogRecordPst(pst),
	}
	encRecord, _ := EncodeLogRecord(record)
	return df.Write(encRecord)
}

func (df *DataFile) Sync() error {
	return df.IoManager.Sync()
}

func (df *DataFile) Close() error {
	return df.IoManager.Close()
}

func (df *DataFile) readNBytes(n int64, offset int64) (b []byte, err error) {
	b = make([]byte, n)
	_, err = df.IoManager.Read(b, offset)
	return
}
