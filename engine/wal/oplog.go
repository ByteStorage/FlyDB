package wal

import (
	"bufio"
	"os"
	"sync"
)

type OpLog struct {
	fileName string
	file     *os.File
	writer   *bufio.Writer
	mu       sync.Mutex
}

func NewOpLog(fileName string) (*OpLog, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &OpLog{
		fileName: fileName,
		file:     file,
		writer:   bufio.NewWriter(file),
	}, nil
}

func (o *OpLog) WriteEntry(entry string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	_, err := o.writer.WriteString(entry + "\n")
	if err != nil {
		return err
	}

	// You can decide when to flush. It can be after every write, or periodically.
	return o.writer.Flush()
}

func (o *OpLog) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if err := o.writer.Flush(); err != nil {
		return err
	}
	return o.file.Close()
}
