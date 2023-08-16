package wal

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"os"
)

type DbWal struct {
	log   *OpLog
	db    *engine.DB
	ch    chan func()
	cache *Cache // add lru cache
	dir   string
}

func NewDbWal(option config.Options, cacheCapacity int) (*DbWal, error) {
	db, err := engine.NewDB(option)
	if err != nil {
		return nil, err
	}

	log, err := NewOpLog(option.DirPath + "/wal")
	if err != nil {
		return nil, err
	}

	d := &DbWal{
		log:   log,
		db:    db,
		ch:    make(chan func(), 1000),
		cache: NewCache(cacheCapacity),
		dir:   option.DirPath,
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
	logMsg := fmt.Sprintf("put key=%s, value=%s", key, value)
	err := d.log.WriteEntry(logMsg)
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
	return nil
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
	}
}
