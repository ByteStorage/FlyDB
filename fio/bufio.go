package fio

import (
	"bufio"
	"os"
)

var _ IOManager = (*Bufio)(nil)

type Bufio struct {
	fd *os.File
	wr *bufio.ReadWriter
}

func (b *Bufio) Read(bytes []byte, offset int64) (int, error) {
	_, err := b.fd.Seek(offset, 0)
	if err != nil {
		return 0, err
	}
	return b.wr.Read(bytes)
}

func (b *Bufio) Write(bytes []byte) (int, error) {
	return b.wr.Write(bytes)
}

func (b *Bufio) Sync() error {
	return b.fd.Sync()
}

func (b *Bufio) Close() error {
	err := b.wr.Flush()
	if err != nil {
		return err
	}
	return b.fd.Close()
}

func (b *Bufio) Size() (int64, error) {
	stat, err := b.fd.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func NewBufIOManager(path string) (*Bufio, error) {
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, DataFilePerm)
	if err != nil {
		return nil, err
	}
	return &Bufio{
		fd: fd,
		wr: bufio.NewReadWriter(bufio.NewReader(fd), bufio.NewWriter(fd)),
	}, nil
}
