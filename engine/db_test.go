package engine

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

func TestNewFlyDB(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb")
	opts.DirPath = dir
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)
}

func TestDB_Put(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-put")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// Put a data file normally
	err = db.Put(randkv.GetTestKey(1), randkv.RandomValue(24))
	assert.Nil(t, err)
	val1, err := db.Get(randkv.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val1)

	// Put the same data repeatedly
	err = db.Put(randkv.GetTestKey(1), randkv.RandomValue(24))
	assert.Nil(t, err)
	val2, err := db.Get(randkv.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val2)

	// The key is empty
	err = db.Put(nil, randkv.RandomValue(24))
	assert.Equal(t, _const.ErrKeyIsEmpty, err)

	// value is null
	err = db.Put(randkv.GetTestKey(22), nil)
	assert.Nil(t, err)
	val3, err := db.Get(randkv.GetTestKey(22))
	assert.Equal(t, 0, len(val3))
	assert.Nil(t, err)

	// The conversion is performed by writing to the data file
	for i := 0; i < 1000000; i++ {
		err := db.Put(randkv.GetTestKey(i), randkv.RandomValue(128))
		assert.Nil(t, err)
	}
	assert.Equal(t, 2, len(db.olderFiles))

	// Restart and Put data again
	err = db.Close()
	assert.Nil(t, err)

	// Restart the database
	db2, err := NewDB(opts)
	defer db2.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db2)
	val4 := randkv.RandomValue(128)
	err = db2.Put(randkv.GetTestKey(55), val4)
	assert.Nil(t, err)
	val5, err := db2.Get(randkv.GetTestKey(55))
	assert.Nil(t, err)
	assert.Equal(t, val4, val5)
}

func TestDB_ConcurrentPut(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-ConcurrentPut")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	var wg sync.WaitGroup

	putTestWorker := func(id int, db *DB) {
		defer func() {
			// fmt.Printf("put Worker %d done\n", id)
			wg.Done()
		}()
		// fmt.Printf("Worker %d processing\n", id)
		err = db.Put(randkv.GetTestKey(id), randkv.GetTestKey(id))
		assert.Nil(t, err)
		// fmt.Printf("Worker %d resumed\n", id)
	}

	getTestWorker := func(id int, db *DB) {
		defer func() {
			wg.Done()
		}()
		val, err := db.Get(randkv.GetTestKey(id))
		assert.Nil(t, err)
		assert.NotNil(t, val)
		assert.Equal(t, randkv.GetTestKey(id), val)
	}

	var workerNum = 100

	// Put workerNum data in parallel
	wg.Add(workerNum)
	go func() {
		for i := 0; i < workerNum; i++ {
			go putTestWorker(i, db)
		}
	}()

	// Parallel Put workerNum to the same data as the previous parallel process
	wg.Add(workerNum)
	go func() {
		for i := 0; i < workerNum; i++ {
			go putTestWorker(i, db)
		}
	}()

	// Wait for all Put to be completed and perform subsequent tests
	wg.Wait()

	// Get all the previously inserted data in parallel
	wg.Add(workerNum)
	go func() {
		for i := 0; i < workerNum; i++ {
			go getTestWorker(i, db)
		}
	}()

	// Wait for the parallel Get to end and conduct subsequent tests
	wg.Wait()

	// Convert to the old data file and get the value from the old data file
	for i := workerNum; i < 1000000; i++ {
		err := db.Put(randkv.GetTestKey(i), randkv.RandomValue(128))
		assert.Nil(t, err)
	}

	wg.Add(workerNum)
	go func() {
		for i := 0; i < workerNum; i++ {
			go getTestWorker(i, db)
		}
	}()
	wg.Wait()

	// After the restart, all the data previously written can be obtained
	err = db.Close()
	assert.Nil(t, err)

	// Restart the database
	db2, err := NewDB(opts)
	defer db2.Clean()

	val1, err := db2.Get(randkv.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val1)
	assert.Equal(t, randkv.GetTestKey(1), val1)

	val2, err := db2.Get(randkv.GetTestKey(3))
	assert.Nil(t, err)
	assert.NotNil(t, val2)
	assert.Equal(t, randkv.GetTestKey(3), val2)

	val3, err := db2.Get(randkv.GetTestKey(5))
	assert.Nil(t, err)
	assert.NotNil(t, val3)
	assert.Equal(t, randkv.GetTestKey(5), val3)

	val4, err := db2.Get(randkv.GetTestKey(99999))
	assert.Nil(t, err)
	assert.NotNil(t, val4)

	// Restart the database and test Get in parallel
	wg.Add(workerNum)
	go func() {
		for i := 0; i < workerNum; i++ {
			go getTestWorker(i, db2)
		}
	}()
	wg.Wait()
	fmt.Printf("All workers done\n")
}

func TestDB_Get(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-get")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// Read a piece of data normally
	err = db.Put(randkv.GetTestKey(11), randkv.RandomValue(24))
	assert.Nil(t, err)
	val1, err := db.Get(randkv.GetTestKey(11))
	assert.Nil(t, err)
	assert.NotNil(t, val1)

	// Read a nonexistent key
	val2, err := db.Get([]byte("some key unknown"))
	assert.Nil(t, val2)
	assert.Equal(t, _const.ErrKeyNotFound, err)

	// The value is read after being repeatedly Put
	err = db.Put(randkv.GetTestKey(22), randkv.RandomValue(24))
	assert.Nil(t, err)
	err = db.Put(randkv.GetTestKey(22), randkv.RandomValue(24))
	val3, err := db.Get(randkv.GetTestKey(22))
	assert.Nil(t, err)
	assert.NotNil(t, val3)

	// Get the value after it is deleted
	err = db.Put(randkv.GetTestKey(33), randkv.RandomValue(24))
	assert.Nil(t, err)
	err = db.Delete(randkv.GetTestKey(33))
	assert.Nil(t, err)
	val4, err := db.Get(randkv.GetTestKey(33))
	assert.Equal(t, 0, len(val4))
	assert.Equal(t, _const.ErrKeyNotFound, err)

	// Convert it to the old data file and obtain the value from the old data file
	for i := 100; i < 1000000; i++ {
		err := db.Put(randkv.GetTestKey(i), randkv.RandomValue(128))
		assert.Nil(t, err)
	}
	assert.Equal(t, 2, len(db.olderFiles))
	val5, err := db.Get(randkv.GetTestKey(101))
	assert.Nil(t, err)
	assert.NotNil(t, val5)

	// After the restart, all the data previously written can be obtained
	err = db.Close()
	assert.Nil(t, err)

	// Restart the database
	db2, err := NewDB(opts)
	defer db2.Clean()
	val6, err := db2.Get(randkv.GetTestKey(11))
	assert.Nil(t, err)
	assert.NotNil(t, val6)
	assert.Equal(t, val1, val6)

	val7, err := db2.Get(randkv.GetTestKey(22))
	assert.Nil(t, err)
	assert.NotNil(t, val7)
	assert.Equal(t, val3, val7)

	val8, err := db2.Get(randkv.GetTestKey(33))
	assert.Equal(t, 0, len(val8))
	assert.Equal(t, _const.ErrKeyNotFound, err)
}

func TestDB_Delete(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-delete")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// Delete an existing key
	err = db.Put(randkv.GetTestKey(11), randkv.RandomValue(128))
	assert.Nil(t, err)
	err = db.Delete(randkv.GetTestKey(11))
	assert.Nil(t, err)
	_, err = db.Get(randkv.GetTestKey(11))
	assert.Equal(t, _const.ErrKeyNotFound, err)

	// Delete a key that does not exist
	err = db.Delete([]byte("unknown key"))
	assert.Nil(t, err)

	// Delete an empty key
	err = db.Delete(nil)
	assert.Equal(t, _const.ErrKeyIsEmpty, err)

	// The value is deleted and Put again
	err = db.Put(randkv.GetTestKey(22), randkv.RandomValue(128))
	assert.Nil(t, err)
	err = db.Delete(randkv.GetTestKey(22))
	assert.Nil(t, err)

	err = db.Put(randkv.GetTestKey(22), randkv.RandomValue(128))
	assert.Nil(t, err)
	val1, err := db.Get(randkv.GetTestKey(22))
	assert.NotNil(t, val1)
	assert.Nil(t, err)

	// After the restart, perform the verification again
	err = db.Close()
	assert.Nil(t, err)

	// Restart the database
	db2, err := NewDB(opts)
	time.Sleep(time.Millisecond * 100)
	defer db2.Clean()
	_, err = db2.Get(randkv.GetTestKey(11))
	assert.Equal(t, _const.ErrKeyNotFound, err)

	val2, err := db2.Get(randkv.GetTestKey(22))
	assert.Nil(t, err)
	assert.Equal(t, val1, val2)
}

func TestDB_GetListKeys(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-ListKey")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// Database is empty
	keys1 := db.GetListKeys()
	assert.Equal(t, 0, len(keys1))

	// Only one piece of data
	err = db.Put(randkv.GetTestKey(10), randkv.GetTestKey(20))
	assert.Nil(t, err)
	keys2 := db.GetListKeys()
	assert.Equal(t, 1, len(keys2))

	err = db.Put(randkv.GetTestKey(20), randkv.GetTestKey(20))
	assert.Nil(t, err)
	err = db.Put(randkv.GetTestKey(30), randkv.GetTestKey(20))
	assert.Nil(t, err)
	err = db.Put(randkv.GetTestKey(40), randkv.GetTestKey(20))
	assert.Nil(t, err)

	keys3 := db.GetListKeys()
	assert.Equal(t, 4, len(keys3))
	for _, value := range keys3 {
		assert.NotNil(t, value)
	}
}

func TestDB_Fold(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-fold")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put(randkv.GetTestKey(10), randkv.RandomValue(20))
	assert.Nil(t, err)
	err = db.Put(randkv.GetTestKey(20), randkv.RandomValue(20))
	assert.Nil(t, err)
	err = db.Put(randkv.GetTestKey(30), randkv.RandomValue(20))
	assert.Nil(t, err)
	err = db.Put(randkv.GetTestKey(40), randkv.RandomValue(20))
	assert.Nil(t, err)

	err = db.Fold(func(key []byte, value []byte) bool {
		assert.NotNil(t, key)
		assert.NotNil(t, value)
		return true
	})
	assert.Nil(t, err)
}

func TestDB_Close(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-close")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put(randkv.GetTestKey(10), randkv.GetTestKey(10))
	assert.Nil(t, err)

	err = db.Close()
	assert.Nil(t, err)
}

func TestDB_Sync(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-close")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put(randkv.GetTestKey(10), randkv.GetTestKey(10))
	assert.Nil(t, err)

	err = db.Sync()
	assert.Nil(t, err)
}
