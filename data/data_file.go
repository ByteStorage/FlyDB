package data

import "flydb/fio"

// DataFile 数据文件
type DataFile struct {
	FileID    uint32        //文件id
	WriteOff  int64         //文件写到了哪个位置
	IoManager fio.IOManager //io 读写操作
}

// OpenDataFile 打开新的数据文件
func OpenDataFile(dirPath string, fildID uint32) (*DataFile, error) {
	return nil, nil
}

func (df *DataFile) ReadLogRecord(offset int64) (*LogRecord, error) {
	return nil, nil
}

func (df *DataFile) Write(buf []byte) error {
	return nil
}

func (df *DataFile) Sync() error {
	return nil
}
