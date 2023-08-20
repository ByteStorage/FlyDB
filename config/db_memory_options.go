package config

import "github.com/ByteStorage/FlyDB/lib/wal"

type DbMemoryOptions struct {
	Option       Options
	LogNum       uint32
	FileSize     int64
	SaveTime     int64
	MemSize      int64
	TotalMemSize int64
	ColumnName   string
	Wal          *wal.Wal
}
