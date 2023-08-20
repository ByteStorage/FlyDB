package memory

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/db/engine"
	"github.com/ByteStorage/FlyDB/lib/wal"
	"os"
	"sync"
)

type Db struct {
	option      config.DbMemoryOptions
	db          *engine.DB
	mem         *MemTable
	oldList     []*MemTable
	wal         *wal.Wal
	oldListChan chan *MemTable
	totalSize   int64
	activeSize  int64
	pool        *sync.Pool
	errMsgCh    chan string
}

// NewDB create a new db of wal and memTable
func NewDB(option config.DbMemoryOptions) (*Db, error) {
	// create a new memTable
	mem := NewMemTable()
	// dir path has been changed to dir path + column name
	option.Option.DirPath = option.Option.DirPath + "/" + option.ColumnName
	db, err := engine.NewDB(option.Option)
	if err != nil {
		return nil, err
	}
	w := option.Wal
	// if wal is nil, create a new wal
	// if wal is not nil, the wal was created by column family
	if option.Wal == nil {
		walOptions := wal.Options{
			DirPath:  option.Option.DirPath,
			FileSize: option.FileSize,
			SaveTime: option.SaveTime,
			LogNum:   option.LogNum,
		}
		w, err = wal.NewWal(walOptions)
		if err != nil {
			return nil, err
		}
	}

	d := &Db{
		mem:         mem,
		db:          db,
		option:      option,
		oldList:     make([]*MemTable, 0),
		oldListChan: make(chan *MemTable, 1000000),
		activeSize:  0,
		totalSize:   0,
		wal:         w,
		pool:        &sync.Pool{New: func() interface{} { return make([]byte, 0, 1024) }},
	}
	go d.async()
	go d.wal.AsyncSave()
	go d.handlerErrMsg()
	return d, nil
}

func (d *Db) handlerErrMsg() {
	log := d.option.Option.DirPath + "/error.log"
	for msg := range d.errMsgCh {
		// write to log
		_ = os.WriteFile(log, []byte(msg), 0666)
	}
}

func (d *Db) Put(key []byte, value []byte) error {
	// calculate key and value size
	keyLen := int64(len(key))
	valueLen := int64(len(value))

	d.pool.Put(func() {
		// Write to wal, try 3 times
		ok := false
		for i := 0; i < 3; i++ {
			err := d.wal.Put(key, value)
			if err == nil {
				ok = true
				break
			}
		}
		if !ok {
			err := d.wal.Delete(key)
			if err != nil {
				d.errMsgCh <- "write to wal error when delete the key: " + string(key) + " error: " + err.Error()
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

func (d *Db) Delete(key []byte) error {
	panic("implement me")
}

func (d *Db) Keys() ([][]byte, error) {
	panic("implement me")
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
