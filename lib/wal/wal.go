package wal

import (
	"github.com/tidwall/wal"
)

const FileName = "/tmp/flydb/wal"

type Wal struct {
	log *wal.Log
}

func (w *Wal) Write(data []byte) error {
	index, err := w.log.LastIndex()
	if err != nil {
		return err
	}
	return w.log.Write(index+1, data)
}

func (w *Wal) Read(index uint64) ([]byte, error) {
	return w.log.Read(index)
}

func (w *Wal) ReadLast() ([]byte, error) {
	index, err := w.log.LastIndex()
	if err != nil {
		return nil, err
	}
	return w.log.Read(index)
}

func New() (*Wal, error) {
	log, err := wal.Open(FileName, nil)
	if err != nil {
		return &Wal{}, err
	}
	index, err := log.LastIndex()
	if index == 0 {
		err := log.Write(1, []byte("--------------------"))
		if err != nil {
			return &Wal{}, err
		}
	}
	return &Wal{log: log}, nil
}
