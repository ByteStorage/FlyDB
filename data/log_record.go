package data

import "encoding/binary"

type LogRecrdType = byte

const (
	LogRecordNormal LogRecrdType = iota
	LogRecordDeleted
)

// crc type KeySize ValueSize
// 4 +  1 +   5   +    5    (byte)
const maxLogRecordHeaderSize = binary.MaxVarintLen32

// LogRecord 写入到数据文件的记录
type LogRecord struct {
	Key   []byte
	Value []byte
	Type  LogRecrdType
}

// LogRecordHeader LogRecord 的头部信息
type LogRecordHeader struct {
	crc        uint32       // crc 校验值
	recordType LogRecrdType // 标识 LogRecord 类型
	keySize    uint32       // key 的长度
	valueSize  uint32       // value 长度
}

// LogRecordPst 数据内存索引， 主要是描述数据在磁盘上的位置
type LogRecordPst struct {
	Fid    uint32 // 文件id，表示将数据存储到了哪个文件当中
	Offset int64  // 偏移，表示将数据存储到了数据文件的哪个位置
}

// EncodeLogRecord 对 LogRecord 进行编码，返回字节数组和长度
func EncodeLogRecord(logrecord *LogRecord) ([]byte, int64) {
	return nil, 0
}

// decodeLogRecordHeader 对字节数组中的 Header 信息进行解码
func decodeLogRecordHeader(buf []byte) (*LogRecordHeader, int64) {
	return nil, 0
}

func getLogRecordCRC(lr *LogRecord, header []byte) uint32 {
	return 0
}
