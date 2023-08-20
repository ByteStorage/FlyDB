package config

import "github.com/ByteStorage/FlyDB/lib/wal"

type ColumnOptions struct {
	DbMemoryOptions DbMemoryOptions
	WalOptions      wal.Options
}
