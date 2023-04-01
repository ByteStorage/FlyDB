package data

type LogRecrdType = byte

const (
	LogRecordNormal LogRecrdType = iota
	LogRecordDeleted
)

// LogRecord 写入到数据文件的记录
type LogRecord struct {
	Key   []byte
	Value []byte
	Type  LogRecrdType
}

//LogRecordPst 数据内存索引， 主要是描述数据在磁盘上的位置
type LogRecordPst struct {
	Fid    uint32 // 文件id，表示将数据存储到了哪个文件当中
	Offset int64  //偏移，表示将数据存储到了数据文件的哪个位置
}

// EncodeLogRecord 对 LogRecord 进行编码，返回字节数组和长度
func EncodeLogRecord(logrecord *LogRecord) ([]byte, int64) {
	return nil, 0
}
