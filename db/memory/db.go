package memory

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/db/engine"
	"strings"
	"sync"
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
	totalSize   int64
	activeSize  int64
	pool        *sync.Pool
	errMsgCh    chan []byte
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
		activeSize:  0,
		totalSize:   0,
		pool:        &sync.Pool{New: func() interface{} { return make([]byte, 0, 1024) }},
	}
	go d.async()
	return d, nil
}

func (d *Db) handlerErrMsg() {
	for msg := range d.errMsgCh {
		// TODO handle error: either log it, retry, or whatever makes sense for your application
		_ = msg
	}
}

func (d *Db) Put(key []byte, value []byte) error {
	// calculate key and value size
	keyLen := int64(len(key))
	valueLen := int64(len(value))

	d.pool.Put(func() {
		// Write to WAL
		err := d.wal.Put(key, value)
		if err != nil {
			err := d.wal.Delete(key)
			if err != nil {
				d.errMsgCh <- []byte(err.Error())
			}
		}
	})

	// if sync write, save wal
	if d.option.Option.SyncWrite {
		err := d.wal.Save()
		if err != nil {
			return err
		}
	}

	// if all memTable size > total memTable size, write to db
	if d.totalSize > d.option.TotalMemSize {
		return d.db.Put(key, value)
	}

	// if active memTable size > define size, change to immutable memTable
	if d.activeSize+keyLen+valueLen > d.option.MemSize {
		// add to immutable memTable list
		d.AddOldMemTable(d.mem)
		// create new active memTable
		d.mem = NewMemTable()
		d.activeSize = 0
	}

	// write to active memTable
	d.mem.Put(string(key), value)

	// add size
	d.activeSize += keyLen + valueLen
	d.totalSize += keyLen + valueLen
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
			d.totalSize -= int64(len(key) + len(value))
		}
	}
}

func (d *Db) Clean() {
	d.db.Clean()
}
