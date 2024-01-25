package config

import (
	"github.com/ByteStorage/FlyDB/lib/wal"
	"os"
	"path/filepath"
)

// Options is a comprehensive configuration struct that
// encapsulates various settings for configuring the behavior of a database.
type Options struct {
	// DirPath specifies the path to the directory where the database will store its data files.
	DirPath string

	// DataFileSize defines the maximum size of each data file in the database.
	DataFileSize int64

	// SyncWrite determines whether the database should ensure data persistence with
	// every write operation.
	SyncWrite bool

	// IndexType selects the type of indexing mechanism to be used for efficient data retrieval.
	IndexType IndexerType

	// FIOType indicates the type of file I/O optimization to be applied by the database.
	FIOType FIOType
}

// ColumnOptions are configurations for database column families
type ColumnOptions struct {
	// DbMemoryOptions contains configuration settings
	// for managing database memory usage and caching.
	DbMemoryOptions DbMemoryOptions

	// WalOptions contains configuration settings for the Write-Ahead Logging (WAL) mechanism.
	WalOptions wal.Options
}

// DbMemoryOptions is related to configuration of database memory tables
type DbMemoryOptions struct {
	// Option contains a set of database configuration options
	// to influence memory management behavior.
	Option Options

	// LogNum specifies the number of logs to keep in memory
	// for efficient access and performance.
	LogNum uint32

	// FileSize defines the maximum size of data files to be kept in memory.
	FileSize int64

	// SaveTime determines the interval at which data should be
	// saved from memory to disk to ensure durability.
	SaveTime int64

	// MemSize sets the limit on the amount of memory to be used for caching purposes.
	MemSize int64

	// TotalMemSize defines the overall memory capacity allocated for database operations.
	TotalMemSize int64

	// ColumnName identifies the specific database column to which these memory options apply.
	ColumnName string

	// Wal is a reference to the Write-Ahead Logging (WAL) mechanism that ensures data durability.
	Wal *wal.Wal
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

	// SkipList
	SkipList

	// ARTWithBloom index With Bloom Filter
	ARTWithBloom
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

var DefaultDbMemoryOptions = DbMemoryOptions{
	Option:       DefaultOptions,
	LogNum:       1000,
	FileSize:     256 * 1024 * 1024, // 256MB
	SaveTime:     100 * 1000,
	MemSize:      256 * 1024 * 1024,      // 256MB
	TotalMemSize: 1 * 1024 * 1024 * 1024, // 2GB
	ColumnName:   "default",
	Wal:          nil,
}

var (
	RedisStringDirPath = filepath.Join(os.TempDir(), "flydb/redis/string")
	RedisHashDirPath   = filepath.Join(os.TempDir(), "flydb/redis/hash")
	RedisListDirPath   = filepath.Join(os.TempDir(), "flydb/redis/list")
	RedisSetDirPath    = filepath.Join(os.TempDir(), "flydb/redis/set")
	RedisZSetDirPath   = filepath.Join(os.TempDir(), "flydb/redis/zset")
)
