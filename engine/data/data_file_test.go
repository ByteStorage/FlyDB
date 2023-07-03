package data

import (
	"github.com/ByteStorage/FlyDB/engine/fileio"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const DefaultFileSize = 256 * 1024 * 1024

func TestOpenDataFile(t *testing.T) {
	dataFile1, err := OpenDataFile(os.TempDir(), 0, DefaultFileSize, fileio.FileIOType)
	assert.Nil(t, err)
	assert.NotNil(t, dataFile1)

	dataFile2, err := OpenDataFile(os.TempDir(), 1, DefaultFileSize, fileio.FileIOType)
	assert.Nil(t, err)
	assert.NotNil(t, dataFile2)

	dataFile3, err := OpenDataFile(os.TempDir(), 1, DefaultFileSize, fileio.FileIOType)
	assert.Nil(t, err)
	assert.NotNil(t, dataFile3)
}

func TestDataFile_Write(t *testing.T) {
	dataFile, err := OpenDataFile(os.TempDir(), 12312, DefaultFileSize, fileio.FileIOType)
	assert.Nil(t, err)
	assert.NotNil(t, dataFile)

	err = dataFile.Write([]byte("abc"))
	assert.Nil(t, err)

	err = dataFile.Write([]byte(" hello"))
	assert.Nil(t, err)

	err = dataFile.Write([]byte(" nihao"))
	assert.Nil(t, err)
}

func TestDataFile_Close(t *testing.T) {
	dataFile, err := OpenDataFile(os.TempDir(), 1111111, DefaultFileSize, fileio.FileIOType)
	assert.Nil(t, err)
	assert.NotNil(t, dataFile)

	err = dataFile.Close()
	assert.Nil(t, err)
}

func TestDataFile_Sync(t *testing.T) {
	dataFile, err := OpenDataFile(os.TempDir(), 2222222, DefaultFileSize, fileio.FileIOType)
	assert.Nil(t, err)
	assert.NotNil(t, dataFile)

	err = dataFile.Sync()
	assert.Nil(t, err)
}

func TestDataFile_ReadLogRecord(t *testing.T) {
	dataFile, err := OpenDataFile(os.TempDir(), 123, DefaultFileSize, fileio.FileIOType)
	assert.Nil(t, err)
	assert.NotNil(t, dataFile)

	// There is only one LogRecord
	record1 := &LogRecord{
		Key:   []byte("name"),
		Value: []byte("flydb"),
		Type:  LogRecordNormal,
	}
	buf1, size := EncodeLogRecord(record1)
	err = dataFile.Write(buf1)
	assert.Nil(t, err)

	readRec1, readSize1, err := dataFile.ReadLogRecord(0)
	assert.Nil(t, err)
	assert.Equal(t, size, readSize1)
	assert.Equal(t, record1, readRec1)

	// Multiple logrecords are read from different locations
	record2 := &LogRecord{
		Key:   []byte("name"),
		Value: []byte("bitcask-kv"),
		Type:  LogRecordNormal,
	}
	buf2, size2 := EncodeLogRecord(record2)
	err = dataFile.Write(buf2)
	assert.Nil(t, err)
	readRec2, readSize2, err := dataFile.ReadLogRecord(16)
	assert.Equal(t, size2, readSize2)
	assert.Equal(t, record2, readRec2)

	// The deleted data is at the end of the data file
	record3 := &LogRecord{
		Key:   []byte("name"),
		Value: []byte("delete-data"),
		Type:  LogRecordNormal,
	}
	buf3, size3 := EncodeLogRecord(record3)
	err = dataFile.Write(buf3)
	assert.Nil(t, err)
	readRec3, readSize3, err := dataFile.ReadLogRecord(size + size2)
	assert.Equal(t, size3, readSize3)
	assert.Equal(t, record3, readRec3)

}
