package wal

// Options encapsulates configuration settings for the Write-Ahead Logging (WAL)
// mechanism in a database.
type Options struct {
	// DirPath specifies the directory path where Write-Ahead Logging (WAL) files will be stored.
	DirPath string

	// FileSize determines the maximum size of individual WAL files.
	FileSize int64

	// SaveTime defines the interval at which WAL data should be persisted from memory to disk.
	SaveTime int64

	// LogNum specifies the number of WAL logs to retain, influencing performance and
	// recovery behavior.
	LogNum uint32
}
