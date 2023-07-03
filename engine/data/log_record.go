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

// TransactionRecord 暂存事务相关的数据
type TransactionRecord struct {
	Record *LogRecord
	Pos    *LogRecordPst
}

// EncodeLogRecord 对 LogRecord 进行编码，返回字节数组和长度
// +-------------+------------+------------+--------------+-------+---------+
// |  crc 校验值  |  type 类型  |  key size  |  value size  |  key  |  value  |
// +-------------+------------+------------+--------------+-------+---------+
// |   4 字节     |  1 字节     | 变长(最大5) |  变长(最大5)   |  变长  |  变长   |
// +-------------+------------+------------+--------------+-------+---------+
func EncodeLogRecord(logrecord *LogRecord) ([]byte, int64) {
	// 初始化一个 header 部分的字节数组
	header := make([]byte, maxLogRecordHeaderSize)

	// 第五个字节存储 type
	header[4] = logrecord.Type
	var headerIndex = 5
	// 5 字节之后存储的是 key 和 value 的长度信息
	// 使用变长变量，节省空间
	headerIndex += binary.PutVarint(header[headerIndex:], int64(len(logrecord.Key)))
	headerIndex += binary.PutVarint(header[headerIndex:], int64(len(logrecord.Value)))

	var size = headerIndex + len(logrecord.Key) + len(logrecord.Value)
	encBytes := make([]byte, size)

	// 拷贝 header 内容和 key/value 数据
	copy(encBytes[:headerIndex], header[:headerIndex])
	copy(encBytes[headerIndex:], logrecord.Key)
	copy(encBytes[headerIndex+len(logrecord.Key):], logrecord.Value)

	// 对整个 LogRecord 数据进行校验
	crc := crc32.ChecksumIEEE(encBytes[4:])
	binary.LittleEndian.PutUint32(encBytes[:4], crc)

	return encBytes, int64(size)
}

// EncodeLogRecordPst 对位置信息进行编码
func EncodeLogRecordPst(pst *LogRecordPst) []byte {
	buf := make([]byte, binary.MaxVarintLen16+binary.MaxVarintLen64)
	var index = 0
	index += binary.PutVarint(buf[index:], int64(pst.Fid))
	index += binary.PutVarint(buf[index:], pst.Offset)
	return buf[:index]
}

// DecodeLogRecordPst 对位置信息进行解码
func DecodeLogRecordPst(buf []byte) *LogRecordPst {
	var index = 0
	fileID, n := binary.Varint(buf[index:])
	index += n
	offset, _ := binary.Varint(buf[index:])
	return &LogRecordPst{
		Fid:    uint32(fileID),
		Offset: offset,
	}
}

// decodeLogRecordHeader 对字节数组中的 Header 信息进行解码
func decodeLogRecordHeader(buf []byte) (*LogRecordHeader, int64) {
	if len(buf) <= 4 {
		return nil, 0
	}

	header := &LogRecordHeader{
		crc:        binary.LittleEndian.Uint32(buf[:4]),
		recordType: buf[4],
	}

	var headerIndex = 5
	// 取出 key/value 对应的 size
	keySize, lens := binary.Varint(buf[headerIndex:])
	header.keySize = uint32(keySize)
	headerIndex += lens

	valueSize, lens := binary.Varint(buf[headerIndex:])
	header.valueSize = uint32(valueSize)
	headerIndex += lens

	return header, int64(headerIndex)
}

// getLogRecordCRC 得到 LoRecord 的 crc 校验值
func getLogRecordCRC(logrecord *LogRecord, header []byte) uint32 {
	if logrecord == nil {
		return 0
	}

	crc := crc32.ChecksumIEEE(header[:])
	crc = crc32.Update(crc, crc32.IEEETable, logrecord.Key)
	crc = crc32.Update(crc, crc32.IEEETable, logrecord.Value)

	return crc
}
