package wal

import "encoding/json"

type WriteMessage struct {
	Key   []byte
	Value []byte
	Type  int //0: put, 1: get , 2: delete , 3: update
}

func (w *Wal) Put(key []byte, value []byte) error {
	marshal, err := json.Marshal(&WriteMessage{
		Key:   key,
		Value: value,
		Type:  0,
	})
	if err != nil {
		return err
	}
	return w.Write(marshal)
}

func (w *Wal) Get(key []byte) error {
	marshal, err := json.Marshal(&WriteMessage{
		Key:   key,
		Type:  1,
		Value: nil,
	})
	if err != nil {
		return err
	}
	return w.Write(marshal)
}

func (w *Wal) Delete(key []byte) error {
	marshal, err := json.Marshal(&WriteMessage{
		Key:   key,
		Type:  2,
		Value: nil,
	})
	if err != nil {
		return err
	}
	return w.Write(marshal)
}

func (w *Wal) Update(key []byte, value []byte) error {
	marshal, err := json.Marshal(&WriteMessage{
		Key:   key,
		Value: value,
		Type:  0,
	})
	if err != nil {
		return err
	}
	return w.Write(marshal)
}
