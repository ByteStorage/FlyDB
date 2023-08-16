package wal

import (
	"bufio"
	"github.com/ByteStorage/FlyDB/engine/fileio"
	"os"
)

type SequentialLogger struct {
	m      *fileio.MMapIO
	writer *bufio.Writer
}

func NewSequentialLogger(filepath string) (*SequentialLogger, error) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	manager, err := fileio.NewMMapIOManager(filepath, 1024*1024*1024)
	return &SequentialLogger{
		m:      manager,
		writer: bufio.NewWriter(f),
	}, nil
}

func (sl *SequentialLogger) Write(data string) error {
	_, err := sl.writer.WriteString(data + "\n")
	return err
}

func (sl *SequentialLogger) Flush() error {
	return sl.writer.Flush()
}
