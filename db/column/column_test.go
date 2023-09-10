package column

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/wal"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var DefaultColumnOptions = config.ColumnOptions{
	DbMemoryOptions: config.DbMemoryOptions{
		Option: config.Options{
			DirPath:      "./",
			DataFileSize: 256 * 1024 * 1024, // 256MB
		},
		LogNum:       1000,
		FileSize:     256 * 1024 * 1024, // 256MB
		SaveTime:     100 * 1000,
		MemSize:      256 * 1024 * 1024,      // 256MB
		TotalMemSize: 1 * 1024 * 1024 * 1024, // 2GB
		ColumnName:   "default",
		Wal:          nil,
	},
	WalOptions: wal.Options{
		DirPath:  "./wal_test",
		LogNum:   100,
		FileSize: 100 * 1024 * 1024,
		SaveTime: 100 * 1000,
	},
}

func CleanWalTest() {
	err := os.RemoveAll("./wal_test")
	if err != nil {
		return
	}
}

func TestColumn_CreateColumnFamily(t *testing.T) {
	option := DefaultColumnOptions
	defer CleanWalTest()

	column, err := NewColumn(option)
	assert.Nil(t, err)
	assert.NotNil(t, column)

	err = column.CreateColumnFamily("test")
	assert.Nil(t, err)

	err = column.CreateColumnFamily("test")
	assert.NotNil(t, err)

	err = column.CreateColumnFamily("test1")
	assert.Nil(t, err)

	err = column.CreateColumnFamily("test2")
	assert.Nil(t, err)

	err = column.DropColumnFamily("test")
	assert.Nil(t, err)

	err = column.DropColumnFamily("test1")
	assert.Nil(t, err)

	err = column.DropColumnFamily("test2")
	assert.Nil(t, err)
}

func TestColumn_ListColumnFamilies(t *testing.T) {
	option := DefaultColumnOptions
	defer CleanWalTest()

	column, err := NewColumn(option)
	assert.Nil(t, err)
	assert.NotNil(t, column)

	err = column.CreateColumnFamily("test")
	assert.Nil(t, err)

	err = column.CreateColumnFamily("test")
	assert.NotNil(t, err)

	err = column.CreateColumnFamily("test1")
	assert.Nil(t, err)

	err = column.CreateColumnFamily("test2")
	assert.Nil(t, err)

	list, err := column.ListColumnFamilies()
	assert.Nil(t, err)
	assert.Equal(t, 4, len(list))

	err = column.DropColumnFamily("test")
	assert.Nil(t, err)

	err = column.DropColumnFamily("test1")
	assert.Nil(t, err)

	err = column.DropColumnFamily("test2")
	assert.Nil(t, err)
}

func TestColumn_Put(t *testing.T) {
	option := DefaultColumnOptions
	defer CleanWalTest()

	column, err := NewColumn(option)
	assert.Nil(t, err)
	assert.NotNil(t, column)

	err = column.CreateColumnFamily("test")
	assert.Nil(t, err)

	err = column.Put("test", []byte("test"), []byte("test"))
	assert.Nil(t, err)

	err = column.Put("test", []byte("test1"), []byte("test1"))
	assert.Nil(t, err)

	err = column.Put("test", []byte("test2"), []byte("test2"))
	assert.Nil(t, err)

	err = column.DropColumnFamily("test")
	assert.Nil(t, err)

}

func TestColumn_Get(t *testing.T) {
	option := DefaultColumnOptions
	defer CleanWalTest()

	column, err := NewColumn(option)
	assert.Nil(t, err)
	assert.NotNil(t, column)

	err = column.CreateColumnFamily("test")
	assert.Nil(t, err)

	err = column.Put("test", []byte("test"), []byte("test"))
	assert.Nil(t, err)

	err = column.Put("test", []byte("test1"), []byte("test1"))
	assert.Nil(t, err)

	err = column.Put("test", []byte("test2"), []byte("test2"))
	assert.Nil(t, err)

	value, err := column.Get("test", []byte("test"))
	assert.Nil(t, err)
	assert.Equal(t, []byte("test"), value)

	value, err = column.Get("test", []byte("test1"))
	assert.Nil(t, err)
	assert.Equal(t, []byte("test1"), value)

	value, err = column.Get("test", []byte("test2"))
	assert.Nil(t, err)
	assert.Equal(t, []byte("test2"), value)

	err = column.DropColumnFamily("test")
	assert.Nil(t, err)

}

func TestColumn_Delete(t *testing.T) {
	option := DefaultColumnOptions
	defer CleanWalTest()

	column, err := NewColumn(option)
	assert.Nil(t, err)
	assert.NotNil(t, column)

	err = column.CreateColumnFamily("test")
	assert.Nil(t, err)

	err = column.Put("test", []byte("test"), []byte("test"))
	assert.Nil(t, err)

	err = column.Put("test", []byte("test1"), []byte("test1"))
	assert.Nil(t, err)

	value, err := column.Get("test", []byte("test"))
	assert.Nil(t, err)
	assert.Equal(t, []byte("test"), value)

	value, err = column.Get("test", []byte("test1"))
	assert.Nil(t, err)
	assert.Equal(t, []byte("test1"), value)

	err = column.Delete("test", []byte("test"))
	assert.Nil(t, err)

	err = column.Delete("test", []byte("test1"))
	assert.Nil(t, err)

	value, err = column.Get("test", []byte("test"))
	assert.NotNil(t, err)
	assert.Nil(t, value)

	value, err = column.Get("test", []byte("test1"))
	assert.NotNil(t, err)
	assert.Nil(t, value)

	err = column.DropColumnFamily("test")
	assert.Nil(t, err)

}

func TestColumn_Keys(t *testing.T) {
	option := DefaultColumnOptions
	defer CleanWalTest()

	column, err := NewColumn(option)
	assert.Nil(t, err)
	assert.NotNil(t, column)

	err = column.CreateColumnFamily("test")
	assert.Nil(t, err)

	err = column.Put("test", []byte("test"), []byte("test"))
	assert.Nil(t, err)

	err = column.Put("test", []byte("test1"), []byte("test1"))
	assert.Nil(t, err)

	err = column.Put("test", []byte("test2"), []byte("test2"))
	assert.Nil(t, err)

	keys, err := column.Keys("test")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(keys))

	err = column.DropColumnFamily("test")
	assert.Nil(t, err)
}
