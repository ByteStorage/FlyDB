package data

/*
LogRecordPst 数据内存索引， 主要是描述数据在磁盘上的位置
*/
type LogRecordPst struct {
	Fid    uint32 // 文件id，表示将数据存储到了哪个文件当中
	Offset int64  //偏移，表示将数据存储到了数据文件的哪个位置
}
