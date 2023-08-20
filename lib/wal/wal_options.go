package wal

type Options struct {
	DirPath  string
	FileSize int64
	SaveTime int64
	LogNum   uint32
}
