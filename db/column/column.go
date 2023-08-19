package column

import (
	"errors"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/db/memory"
	"github.com/ByteStorage/FlyDB/lib/wal"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// Column is a column family
type Column interface {
	// CreateColumnFamily create column family
	CreateColumnFamily(name string) error
	// DropColumnFamily drop column family
	DropColumnFamily(name string) error
	// ListColumnFamilies list column families
	ListColumnFamilies() ([]string, error)
	// Put a key/value pair into the column family
	Put(cf string, key []byte, value []byte) error
	// Get a value from the column family
	Get(cf string, key []byte) ([]byte, error)
	// Delete a key from the column family
	Delete(cf string, key []byte) error
	// Keys returns all keys in the column family
	Keys(cf string) ([][]byte, error)
}

// NewColumn create a column family
func NewColumn(option config.ColumnOptions) (Column, error) {
	// create wal, all column family share a wal
	walConf := config.WalConfig{
		DirPath:  option.Option.DirPath,
		LogNum:   option.LogNum,
		SaveTime: option.SaveTime,
		FileSize: option.FileSize,
	}
	w, err := wal.NewWal(walConf)
	if err != nil {
		return nil, err
	}

	// load column family
	col, err := loadColumn(option)
	if err != nil {
		return nil, err
	}

	// if column family exists, return it
	if len(col) > 0 {
		columnFamily := make(map[string]*memory.Db)
		for k, v := range col {
			columnFamily[k] = v
		}
		return &column{
			option:       option,
			mux:          sync.RWMutex{},
			columnFamily: columnFamily,
			wal:          w,
		}, nil
	}

	// if column family not exists, create a new column family
	if option.ColumnName == "" {
		option.ColumnName = "default"
	}

	options := config.DbMemoryOptions{
		Option:       option.Option,
		LogNum:       option.LogNum,
		FileSize:     option.FileSize,
		SaveTime:     option.SaveTime,
		MemSize:      option.MemSize,
		TotalMemSize: option.TotalMemSize,
		ColumnName:   option.ColumnName,
	}

	// create a new db
	db, err := memory.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &column{
		option: option,
		mux:    sync.RWMutex{},
		columnFamily: map[string]*memory.Db{
			option.ColumnName: db,
		},
		wal: w,
	}, nil
}

type column struct {
	mux          sync.RWMutex
	wal          *wal.Wal
	columnFamily map[string]*memory.Db
	option       config.ColumnOptions
}

func (c *column) CreateColumnFamily(name string) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	if _, ok := c.columnFamily[name]; ok {
		return errors.New("column family already exists")
	}
	options := config.DbMemoryOptions{
		Option:       c.option.Option,
		LogNum:       c.option.LogNum,
		FileSize:     c.option.FileSize,
		SaveTime:     c.option.SaveTime,
		MemSize:      c.option.MemSize,
		TotalMemSize: c.option.TotalMemSize,
		ColumnName:   c.option.ColumnName,
	}
	db, err := memory.NewDB(options)
	if err != nil {
		return err
	}
	c.columnFamily[name] = db
	return nil
}

func (c *column) DropColumnFamily(name string) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	if _, ok := c.columnFamily[name]; !ok {
		return errors.New("column family not exists")
	}
	delete(c.columnFamily, name)
	return nil
}

func (c *column) ListColumnFamilies() ([]string, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	var list []string
	for k := range c.columnFamily {
		list = append(list, k)
	}
	return list, nil
}

func (c *column) Put(cf string, key []byte, value []byte) error {
	return c.columnFamily[cf].Put(key, value)
}

func (c *column) Get(cf string, key []byte) ([]byte, error) {
	return c.columnFamily[cf].Get(key)
}

func (c *column) Delete(cf string, key []byte) error {
	return c.columnFamily[cf].Delete(key)
}

func (c *column) Keys(cf string) ([][]byte, error) {
	return c.columnFamily[cf].Keys()
}

func loadColumn(option config.ColumnOptions) (map[string]*memory.Db, error) {
	base := option.Option.DirPath
	base = strings.Trim(base, "/")
	// Check if the base path exists
	if _, err := os.Stat(base); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", base)
	}
	// List all directories under the base path
	dirs, err := ioutil.ReadDir(base)
	if err != nil {
		return nil, err
	}
	columns := make(map[string]*memory.Db)
	for _, dir := range dirs {
		if dir.IsDir() {
			colName := dir.Name()
			dirPath := base + "/" + colName
			option.Option.DirPath = dirPath
			options := config.DbMemoryOptions{
				Option:       option.Option,
				LogNum:       option.LogNum,
				FileSize:     option.FileSize,
				SaveTime:     option.SaveTime,
				MemSize:      option.MemSize,
				TotalMemSize: option.TotalMemSize,
				ColumnName:   option.ColumnName,
			}
			db, err := memory.NewDB(options)
			if err != nil {
				return nil, err
			}
			columns[colName] = db
		}
	}
	return columns, nil
}
