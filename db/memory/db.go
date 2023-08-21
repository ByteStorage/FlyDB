package memory

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/db/engine"
	"github.com/ByteStorage/FlyDB/lib/wal"
	"io"
	"log"
	"os"
	"sync"
)

const (
	// Record types
	putType    = byte(1)
	deleteType = byte(2)
)

type Db struct {
	option            config.DbMemoryOptions
	db                *engine.DB
	mem               *MemTable
	oldList           []*MemTable
	wal               *wal.Wal
	oldListChan       chan *MemTable
	totalSize         int64
	activeSize        int64
	pool              *sync.Pool
	errMsgCh          chan string
	mux               sync.RWMutex
	walDataMtList     []*MemTable
	walDataMtListChan chan *MemTable
}

// NewDB create a new db of wal and memTable
func NewDB(option config.DbMemoryOptions) (*Db, error) {
	// create a new memTable
	mem := NewMemTable()

	// dir path has been changed to dir path + column name
	option.Option.DirPath = option.Option.DirPath + "/" + option.ColumnName
	option.Option.IndexType = config.ARTWithBloom
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

	// initialize db
	d := &Db{
		mem:               mem,
		db:                db,
		option:            option,
		oldList:           make([]*MemTable, 0),
		oldListChan:       make(chan *MemTable, 1000000),
		activeSize:        0,
		totalSize:         0,
		wal:               w,
		pool:              &sync.Pool{New: func() interface{} { return make([]byte, 0, 1024) }},
		mux:               sync.RWMutex{},
		walDataMtList:     make([]*MemTable, 0),
		walDataMtListChan: make(chan *MemTable, 1000000),
	}

	// when loading, the system will execute the every record in wal
	d.load()
	// async write to db
	go d.async()
	// async save wal
	go d.wal.AsyncSave()
	// async handler error message
	go d.handlerErrMsg()
	return d, nil
}

func (d *Db) handlerErrMsg() {
	msgLog := d.option.Option.DirPath + "/error.log"
	for msg := range d.errMsgCh {
		// write to log
		_ = os.WriteFile(msgLog, []byte(msg), 0666)
	}
}

var putTypeInt = int64(1)

func (d *Db) Put(key []byte, value []byte) error {
	d.mux.Lock()
	defer d.mux.Unlock()
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
		if putTypeInt == 1 {
			d.addOldMemTable(d.mem)
		} else {
			d.addWalDataToMemTable(d.mem)
			putTypeInt = 1
		}
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
	d.mux.RLock()
	defer d.mux.RUnlock()
	// first get from memTable
	value, err := d.mem.Get(string(key))
	if err == nil {
		return value, nil
	}

	mtList := append(append([]*MemTable(nil), d.walDataMtList...), d.oldList...)

	// if active memTable not found, get from immutable memTable
	for _, list := range mtList {
		value, err = list.Get(string(key))
		if err == nil {
			return value, nil
		}
	}

	// if immutable memTable not found, get from db
	return d.db.Get(key)
}

func (d *Db) Delete(key []byte) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	d.pool.Put(func() {
		// Write to wal, try 3 times
		ok := false
		for i := 0; i < 3; i++ {
			err := d.wal.Delete(key)
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
	// get from active memTable
	get, err := d.mem.Get(string(key))
	if err == nil {
		d.activeSize -= int64(len(key) + len(get))
		d.totalSize -= int64(len(key) + len(get))
		d.mem.Delete(string(key))
		return nil
	}
	// get from immutable memTable
	mtList := append(append([]*MemTable(nil), d.walDataMtList...), d.oldList...)
	for _, list := range mtList {
		get, err = list.Get(string(key))
		if err == nil {
			d.totalSize -= int64(len(key) + len(get))
			list.Delete(string(key))
			return nil
		}
	}
	// get from db
	return d.db.Delete(key)
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

func (d *Db) addOldMemTable(oldList *MemTable) {
	d.oldListChan <- oldList
}

func (d *Db) addWalDataToMemTable(walDataMt *MemTable) {
	d.walDataMtListChan <- walDataMt
}

func (d *Db) async() {
	for oldList := range d.oldListChan {
		for key, value := range oldList.table {
			// Write to db, try 3 times
			ok := false
			for i := 0; i < 3; i++ {
				err := d.db.Put([]byte(key), value)
				if err == nil {
					ok = true
					break
				}
			}
			if !ok {
				err := d.wal.Delete([]byte(key))
				if err != nil {
					d.errMsgCh <- "write to wal error when delete the key: " + string(key) + " error: " + err.Error()
				}
			}
			d.totalSize -= int64(len(key) + len(value))
		}
	}
}

func (d *Db) Clean() {
	d.db.Clean()
}

func (d *Db) load() {
	// Initialize reading from the start of the WAL.
	d.wal.InitReading()

	for {
		record, err := d.wal.ReadNext()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Handle the error: log it, panic, return, etc.
			log.Printf("Error reading from WAL: %v", err)
			return
		}

		switch record.Type {
		case putType:
			// Assuming Db has a Put method
			putTypeInt = 0
			err := d.Put(record.Key, record.Value)
			if err != nil {
				// Handle the error: log it, panic, return, etc.
				log.Printf("Error applying PUT from WAL: %v", err)
			}
		case deleteType:
			// Assuming Db has a Delete method
			err := d.Delete(record.Key)
			if err != nil {
				// Handle the error: log it, panic, return, etc.
				log.Printf("Error applying DELETE from WAL: %v", err)
			}
		default:
			// Handle unknown type.
			log.Printf("Unknown record type in WAL: %v", record.Type)
		}
	}
}
