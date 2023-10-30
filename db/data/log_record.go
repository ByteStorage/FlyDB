package data

import (
	"encoding/binary"
	"hash/crc32"
)

type LogRecrdType = byte

const (
	LogRecordNormal LogRecrdType = iota
	LogRecordDeleted
	LogRecordTransFinished
)

// crc type KeySize ValueSize
// 4 +  1 +   5   +    5    (byte)
const maxLogRecordHeaderSize = binary.MaxVarintLen32*2 + 5

// LogRecord represents a record written to the data file.
type LogRecord struct {
	Key   []byte       // The key of the record
	Value []byte       // The value of the record
	Type  LogRecrdType // The type of the record
}

// LogRecordHeader represents the header information of a LogRecord.
type LogRecordHeader struct {
	crc        uint32       // CRC checksum value
	recordType LogRecrdType // Identifies the type of LogRecord
	keySize    uint32       // Length of the key
	valueSize  uint32       // Length of the value
}

// LogRecordPst represents the in-memory index of data,
// mainly describing the location of data on disk.
type LogRecordPst struct {
	Fid    uint32 // File ID: Indicates which file the data is stored in
	Offset int64  // Offset: Indicates the position in the data file where the data is stored
}

// TransactionRecord temporarily holds transaction-related data.
type TransactionRecord struct {
	Record *LogRecord    // The log record associated with the transaction
	Pos    *LogRecordPst // The position of the log record in the data file
}

// EncodeLogRecord encodes a LogRecord and returns the byte array and length.
// +-------------+------------+------------+--------------+-------+---------+---------+---------+---------+
// |  crc checksum |  record type  |     key size      |      value size     |     key    |  value    |
// +-------------+------------+------------+--------------+-------+---------+
// |   4 bytes     |  1 byte       | variable (max 5)  |  variable (max 5)   |  variable  |  variable |
// +-------------+------------+------------+--------------+-------+---------+---------+---------+---------+
func EncodeLogRecord(logrecord *LogRecord) ([]byte, int64) {
	header := make([]byte, maxLogRecordHeaderSize)

	// Store the record type at the fifth byte
	header[4] = logrecord.Type
	var headerIndex = 5

	// Store the lengths of key and value after the fifth byte
	// Use variable-length encoding to save space
	headerIndex += binary.PutVarint(header[headerIndex:], int64(len(logrecord.Key)))
	headerIndex += binary.PutVarint(header[headerIndex:], int64(len(logrecord.Value)))

	var size = headerIndex + len(logrecord.Key) + len(logrecord.Value)
	encBytes := make([]byte, size)

	// Copy the header content and key/value data
	copy(encBytes[:headerIndex], header[:headerIndex])
	copy(encBytes[headerIndex:], logrecord.Key)
	copy(encBytes[headerIndex+len(logrecord.Key):], logrecord.Value)

	// Calculate the CRC checksum for the entire LogRecord data
	crc := crc32.ChecksumIEEE(encBytes[4:])
	binary.LittleEndian.PutUint32(encBytes[:4], crc)

	return encBytes, int64(size)
}

// EncodeLogRecordPst encodes the position information of a log record.
func EncodeLogRecordPst(pst *LogRecordPst) []byte {
	buf := make([]byte, binary.MaxVarintLen16+binary.MaxVarintLen64)
	var index = 0
	index += binary.PutVarint(buf[index:], int64(pst.Fid)) // Encode file ID
	index += binary.PutVarint(buf[index:], pst.Offset)     // Encode offset
	return buf[:index]
}

// DecodeLogRecordPst decodes the position information from a byte buffer.
func DecodeLogRecordPst(buf []byte) *LogRecordPst {
	var index = 0
	fileID, n := binary.Varint(buf[index:]) // Decode file ID
	index += n
	offset, _ := binary.Varint(buf[index:]) // Decode offset
	return &LogRecordPst{
		Fid:    uint32(fileID), // Convert file ID to uint32
		Offset: offset,         // Assign offset
	}
}

// decodeLogRecordHeader decodes the header information from a byte buffer.
func decodeLogRecordHeader(buf []byte) (*LogRecordHeader, int64) {
	if len(buf) <= 4 {
		return nil, 0
	}

	header := &LogRecordHeader{
		crc:        binary.LittleEndian.Uint32(buf[:4]), // Decode CRC checksum
		recordType: buf[4],                              // Decode record type
	}

	var headerIndex = 5
	// Extract key/value sizes
	keySize, lens := binary.Varint(buf[headerIndex:]) // Decode key size
	header.keySize = uint32(keySize)
	headerIndex += lens

	valueSize, lens := binary.Varint(buf[headerIndex:]) // Decode value size
	header.valueSize = uint32(valueSize)
	headerIndex += lens

	return header, int64(headerIndex)
}

// getLogRecordCRC calculates the CRC checksum for a LogRecord.
func getLogRecordCRC(logRecord *LogRecord, header []byte) uint32 {
	if logRecord == nil {
		return 0
	}
	// Calculate CRC checksum for the header
	crc := crc32.ChecksumIEEE(header[:])
	// Update CRC checksum with the key
	crc = crc32.Update(crc, crc32.IEEETable, logRecord.Key)
	// Update CRC checksum with the value
	crc = crc32.Update(crc, crc32.IEEETable, logRecord.Value)

	return crc
}
