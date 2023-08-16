package wal

import (
	"bufio"
	"os"
)

type SequentialLogger struct {
	file   *os.File
	writer *bufio.Writer
}

func NewSequentialLogger(filepath string) (*SequentialLogger, error) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &SequentialLogger{
		file:   f,
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
