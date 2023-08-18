package memory

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/db/engine"
	"strings"
)

const (
	walFileName = "/db.wal"
)

type Options struct {
	Option       config.Options
	LogNum       uint32
	FileSize     int64
	SaveTime     int64
	MemSize      int64
	TotalMemSize int64
}

type Db struct {
	option      Options
	db          *engine.DB
	mem         *MemTable
	oldList     []*MemTable
	wal         *Wal
	oldListChan chan *MemTable
	size        int64
}

func NewDB(option Options) (*Db, error) {
	mem := NewMemTable()
	db, err := engine.NewDB(option.Option)
	if err != nil {
		return nil, err
	}
	// Create or open the WAL file.
	option.Option.DirPath = strings.TrimSuffix(option.Option.DirPath, "/")
	wal, err := NewWal(option)
	if err != nil {
		return nil, err
	}
	go wal.AsyncSave()
	d := &Db{
		mem:         mem,
		db:          db,
		wal:         wal,
		option:      option,
		oldList:     make([]*MemTable, 0),
		oldListChan: make(chan *MemTable, 1000000),
	}
	go d.async()
	return d, nil
}

func (d *Db) Put(key []byte, value []byte) error {
	// Write to WAL
	err := d.wal.Put(key, value)
	if err != nil {
		return err
	}

	// if sync write, save wal
	if d.option.Option.SyncWrite {
		err := d.wal.Save()
		if err != nil {
			return err
		}
	}

	// if all memTable size > total memTable size, write to db
	if d.size > d.option.TotalMemSize {
		return d.db.Put(key, value)
	}

	// if active memTable size > define size, change to immutable memTable
	if d.mem.Size()+int64(len(key)+len(value)) > d.option.MemSize {
		// add to immutable memTable list
		d.AddOldMemTable(d.mem)
		// add to size
		d.size += d.mem.Size()
		// create new active memTable
		d.mem = NewMemTable()
	}

	// write to active memTable
	d.mem.Put(string(key), value)
	return nil
}

func (d *Db) Get(key []byte) ([]byte, error) {
	// first get from memTable
	value, err := d.mem.Get(string(key))
	if err == nil {
		return value, nil
	}

	// if active memTable not found, get from immutable memTable
	for _, list := range d.oldList {
		value, err = list.Get(string(key))
		if err == nil {
			return value, nil
		}
	}

	// if immutable memTable not found, get from db
	return d.db.Get(key)
}

func (d *Db) Close() error {
	err := d.wal.Save()
	if err != nil {
		return err
	}
	return d.db.Close()
}

func (d *Db) AddOldMemTable(oldList *MemTable) {
	d.oldListChan <- oldList
}

func (d *Db) async() {
	for oldList := range d.oldListChan {
		for key, value := range oldList.table {
			err := d.db.Put([]byte(key), value)
			if err != nil {
				// TODO handle error: either log it, retry, or whatever makes sense for your application
			}
		}
	}
}
