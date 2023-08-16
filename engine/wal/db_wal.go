package wal

import (
	"encoding/binary"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"os"
	"sync"
)

type DbWal struct {
	log   *SequentialLogger
	db    *engine.DB
	ch    chan func()
	cache *Cache // add lru cache
	dir   string
	w     *sync.WaitGroup
}

func NewDbWal(option config.Options, cacheCapacity int) (*DbWal, error) {
	db, err := engine.NewDB(option)
	if err != nil {
		return nil, err
	}

	log, err := NewSequentialLogger(option.DirPath + "/log")
	if err != nil {
		return nil, err
	}

	d := &DbWal{
		log:   log,
		db:    db,
		ch:    make(chan func(), 1000000),
		cache: NewCache(cacheCapacity),
		dir:   option.DirPath,
		w:     &sync.WaitGroup{},
	}

	go d.asyncWorker()
	return d, nil
}

func (d *DbWal) asyncWorker() {
	for task := range d.ch {
		task()
	}
}

func (d *DbWal) PutByWal(key []byte, value []byte) error {
	keySize := len(key)
	valueSize := len(value)
	logMsg := make([]byte, 4+4+keySize+valueSize) // Assuming 4 bytes for each size
	binary.BigEndian.PutUint32(logMsg[0:4], uint32(keySize))
	binary.BigEndian.PutUint32(logMsg[4:8], uint32(valueSize))
	copy(logMsg[8:8+keySize], key)
	copy(logMsg[8+keySize:], value)
	err := d.log.Write(string(logMsg))
	if err != nil {
		return err
	}
	// update cache
	d.cache.Put(key, value)

	// sync update db
	d.ch <- func() {
		if err := d.db.Put(key, value); err != nil { // Add error handling
			fmt.Printf("Error updating DB: %v\n", err) // Log the error
		}
	}
	return err
}

func (d *DbWal) GetByWal(key []byte) ([]byte, error) {
	// get from cache
	value, found := d.cache.Get(key)
	if found {
		return value, nil
	}

	// if not found in cache, get from disk
	valueFromDisk, err := d.db.Get(key)
	if err != nil {
		return nil, err
	}

	// update cache
	d.cache.Put(key, valueFromDisk)

	return valueFromDisk, nil
}

func (d *DbWal) Clean() {
	if d.db != nil {
		_ = d.db.Close()
		err := os.RemoveAll(d.dir)
		if err != nil {
			panic(err)
		}
		err = d.log.Flush()
		if err != nil {
			panic(err)
		}
	}
}
