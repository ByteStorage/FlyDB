package wal

import (
	"github.com/ByteStorage/FlyDB/engine/fileio"
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

func (sl *SequentialLogger) Write(data string) error {
	_, err := sl.m.Write([]byte(data + "\n"))
	return err
}

func (sl *SequentialLogger) Flush() error {
	return sl.m.Sync()
}
