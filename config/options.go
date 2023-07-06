package config

import "os"

type Options struct {
	DirPath      string // Database data directory
	DataFileSize int64  // Size of data files
	SyncWrite    bool   // Whether to persist data on every write
	IndexType    IndexerType
	FIOType      FIOType
}

// IteratorOptions is the configuration for index iteration.
type IteratorOptions struct {
	// Prefix specifies the prefix value for keys to iterate over. Default is empty.
	Prefix []byte
	// Reverse indicates whether to iterate in reverse order.
	// Default is false for forward iteration.
	Reverse bool
}

// WriteBatchOptions is the configuration for batch writing.
type WriteBatchOptions struct {
	// MaxBatchNum is the maximum number of data entries in a batch.
	MaxBatchNum uint
	// SyncWrites indicates whether to sync (persist) the data on batch commit.
	SyncWrites bool
}

type FIOType = int8

const (
	FileIOType = iota + 1 // Standard File IO
	BufIOType             // File IO with buffer
	MmapIOType            // Memory Mapping IO
)

type IndexerType = int8

const (
	// Btree
	Btree IndexerType = iota + 1

	// ART (Adpative Radix Tree)
	ART
)

const (
	DefaultAddr = "127.0.0.1:8999"
)

var DefaultOptions = Options{
	DirPath:      os.TempDir(),
	DataFileSize: 256 * 1024 * 1024, // 256MB
	SyncWrite:    false,
	IndexType:    ART,
	FIOType:      MmapIOType,
}

var DefaultIteratorOptions = IteratorOptions{
	Prefix:  nil,
	Reverse: false,
}

var DefaultWriteBatchOptions = WriteBatchOptions{
	MaxBatchNum: 10000,
	SyncWrites:  true,
}
