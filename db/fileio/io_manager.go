package fileio

// 0644 Indicates that a file is created.
// The file owner can read and write the file,
// but others can only read the file
const DataFilePerm = 0644

const DefaultFileSize = 256 * 1024 * 1024

const (
	FileIOType = iota + 1 // Standard File IO
	BufIOType             // File IO with buffer
	MmapIOType            // Memory Mapping IO
)

// IOManager is an abstract IO management interface that can accommodate different IO types.
// Currently, it supports standard file IO.
type IOManager interface {
	// Read reads the corresponding data from the file at the given position.
	Read([]byte, int64) (int, error)

	// Write writes a byte array to the file.
	Write([]byte) (int, error)

	// Sync persists data.
	Sync() error

	// Close closes the file.
	Close() error

	// Size gets the file size.
	Size() (int64, error)
}

// NewIOManager get IOManager based on type
func NewIOManager(filename string, fileSize int64, fioType int8) (IOManager, error) {
	switch fioType {
	case FileIOType:
		return NewFileIOManager(filename)
	case BufIOType:
		return NewBufIOManager(filename)
	case MmapIOType:
		return NewMMapIOManager(filename, fileSize)
	}
	return NewMMapIOManager(filename, fileSize)
}
