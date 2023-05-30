package fio

const DataFilePerm = 0644 //0644 表示创建了一个文件，文件所有者可以读写，其他人只能读

// IOManager 抽象 IO 管理接口， 可以接入不同的 IO 类型， 目前支持标准文件 IO
type IOManager interface {
	// Read 从文件的给定位置读取对应的数据
	Read([]byte, int64) (int, error)

	// Write 写入字节数组到文件中
	Write([]byte) (int, error)

	// Sync 持久化数据
	Sync() error

	// Close 关闭文件
	Close() error

	// Size get file size
	Size() (int64, error)
}
