package data

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/fio"
	"hash/crc32"
	"io"
	"path/filepath"
)

const (
	DataFileSuffix      = ".data"
	HintFileSuffix      = "hintIndex"
	MergeFinaFileSuffix = "mergeFina"
)

// DataFile 数据文件
type DataFile struct {
	FileID    uint32        //文件id
	WriteOff  int64         //文件写到了哪个位置
	IoManager fio.IOManager //io 读写操作
}

// OpenDataFile 打开新的数据文件
func OpenDataFile(dirPath string, fildID uint32, fileSize int64, fioType int8) (*DataFile, error) {
	fileName := GetDataFileName(dirPath, fildID)
	return newDataFile(fileName, fildID, fileSize, fioType)
}

func GetDataFileName(dirPath string, fildID uint32) string {
	return filepath.Join(dirPath, fmt.Sprintf("%09d", fildID)+DataFileSuffix)
}

// OpenHintFile 打开 Hint 索引文件
func OpenHintFile(dirPath string, fileSize int64, fioType int8) (*DataFile, error) {
	fileName := filepath.Join(dirPath, HintFileSuffix)
	return newDataFile(fileName, 0, fileSize, fioType)
}

// OpenMergeFinaFile 打开标识 merge 完成的文件
func OpenMergeFinaFile(dirPath string, fileSize int64, fioType int8) (*DataFile, error) {
	fileName := filepath.Join(dirPath, MergeFinaFileSuffix)
	return newDataFile(fileName, 0, fileSize, fioType)
}

func newDataFile(dirPath string, fildID uint32, fileSize int64, fioType int8) (*DataFile, error) {
	//初始化 IOManager 管理器接口
	ioManager, err := fio.NewIOManager(dirPath, fileSize, fioType)
	if err != nil {
		return nil, err
	}
	return &DataFile{
		FileID:    fildID,
		WriteOff:  0,
		IoManager: ioManager,
	}, nil
}

// ReadLogRecord 根据 offset 从数据文件中读取 logRecord
func (df *DataFile) ReadLogRecord(offset int64) (*LogRecord, int64, error) {
	fileSize, err := df.IoManager.Size()
	if err != nil {
		return nil, 0, err
	}

	var headerBytes int64 = maxLogRecordHeaderSize
	if offset+maxLogRecordHeaderSize > fileSize {
		headerBytes = fileSize - offset
	}

	// 读取 header 信息
	headerBuf, err := df.readNBytes(headerBytes, offset)
	if err != nil {
		return nil, 0, err
	}

	header, headerSize := decodeLogRecordHeader(headerBuf)
	// 下面俩个条件表示读到了文件末尾，直接返回 EOF
	if header == nil {
		return nil, 0, io.EOF
	}
	if header.crc == 0 && header.keySize == 0 && header.valueSize == 0 {
		return nil, 0, io.EOF
	}

	// 取出对应的 key 和 value 的长度
	keySize, valueSize := int64(header.keySize), int64(header.valueSize)
	var recordSize = headerSize + keySize + valueSize

	logRecord := &LogRecord{Type: header.recordType}

	// 读取用户实际存储的 key/value 数据
	if keySize > 0 || valueSize > 0 {
		kvBuf, err := df.readNBytes(keySize+valueSize, headerSize+offset)
		if err != nil {
			return nil, 0, err
		}
		// 解码
		logRecord.Key = kvBuf[:keySize]
		logRecord.Value = kvBuf[keySize:]
	}

	// 校验 crc （检查数据的有效性）
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

// WriteHintRecord 写入索引信息到 hint 文件中
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
