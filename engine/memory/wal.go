package memory

import (
	"github.com/ByteStorage/FlyDB/engine/fileio"
)

const (
	// TypeData Record types
	TypeData byte = iota // Just for demonstration, you can define more types if needed
)

type SequentialLogger struct {
	m *fileio.MMapIO
}

func NewSequentialLogger(filepath string) (*SequentialLogger, error) {
	manager, err := fileio.NewMMapIOManager(filepath, 1024*1024*1024)
	if err != nil {
		return nil, err
	}
	return &SequentialLogger{
		m: manager,
	}, nil
}

// Put writes a record to the log.
//+---------+-----------+-----------+----------------+--- ... ---+
//|CRC (4B) | Size (2B) | Type (1B) | Log number (4B)| Payload   |
//+---------+-----------+-----------+----------------+--- ... ---+
//Same as above, with the addition of
//Log number = 32bit log file number, so that we can distinguish between
//records written by the most recent log writer vs a previous one.
func (sl *SequentialLogger) Put(key []byte, value []byte) error {
	panic("implement me")
}

func (sl *SequentialLogger) Flush() error {
	return sl.m.Sync()
}
