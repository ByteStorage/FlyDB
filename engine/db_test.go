package engine

import (
	"fmt"
	"github.com/ByteStorage/flydb"
	"github.com/ByteStorage/flydb/config"
	"github.com/ByteStorage/flydb/lib/randkv"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

// 测试完成之后销毁 DB 数据目录
func destroyDB(db *DB) {
	if db != nil {
		if db.activeFile != nil {
			_ = db.Close()
		}
		err := os.RemoveAll(db.options.DirPath)
		if err != nil {
			panic(err)
		}
	}
}

func TestNewFlyDB(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb")
	opts.DirPath = dir
	db, err := NewDB(opts)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)
}

func TestDB_Put(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-put")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// 1.正常 Put 一条数据
	err = db.Put(randkv.GetTestKey(1), randkv.RandomValue(24))
	assert.Nil(t, err)
	val1, err := db.Get(randkv.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val1)

	// 2.重复 Put key 相同的数据
	err = db.Put(randkv.GetTestKey(1), randkv.RandomValue(24))
	assert.Nil(t, err)
	val2, err := db.Get(randkv.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val2)

	// 3.key 为空
	err = db.Put(nil, randkv.RandomValue(24))
	assert.Equal(t, flydb.ErrKeyIsEmpty, err)

	// 4.value 为空
	err = db.Put(randkv.GetTestKey(22), nil)
	assert.Nil(t, err)
	val3, err := db.Get(randkv.GetTestKey(22))
	assert.Equal(t, 0, len(val3))
	assert.Nil(t, err)

	// 5.写到数据文件进行了转换
	for i := 0; i < 1000000; i++ {
		err := db.Put(randkv.GetTestKey(i), randkv.RandomValue(128))
		assert.Nil(t, err)
	}
	assert.Equal(t, 2, len(db.olderFiles))

	// 6.重启后再 Put 数据
	err = db.Close()
	assert.Nil(t, err)

	// 重启数据库
	db2, err := NewDB(opts)
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
	defer destroyDB(db)
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

	// 1. 并行Put  workerNum条数据
	wg.Add(workerNum)
	go func() {
		for i := 0; i < workerNum; i++ {
			go putTestWorker(i, db)
		}
	}()

	// 2. 并行Put  workerNum条和上一并行过程相同的数据
	wg.Add(workerNum)
	go func() {
		for i := 0; i < workerNum; i++ {
			go putTestWorker(i, db)
		}
	}()

	// 等待所有Put结束后，进行后续测试
	wg.Wait()

	// 3. 并行Get 之前插入的所有数据
	wg.Add(workerNum)
	go func() {
		for i := 0; i < workerNum; i++ {
			go getTestWorker(i, db)
		}
	}()

	// 等待并行Get结束后，进行后续测试
	wg.Wait()

	// 4. 转换为了旧的数据文件，从旧的数据文件上获取 value
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

	// 6.重启后，前面写入的数据都能拿到
	err = db.Close()
	assert.Nil(t, err)

	// 重启数据库
	db2, err := NewDB(opts)

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

	// 重启数据库后再并行测试Get
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
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// 1.正常读取一条数据
	err = db.Put(randkv.GetTestKey(11), randkv.RandomValue(24))
	assert.Nil(t, err)
	val1, err := db.Get(randkv.GetTestKey(11))
	assert.Nil(t, err)
	assert.NotNil(t, val1)

	// 2.读取一个不存在的 key
	val2, err := db.Get([]byte("some key unknown"))
	assert.Nil(t, val2)
	assert.Equal(t, flydb.ErrKeyNotFound, err)

	// 3.值被重复 Put 后在读取
	err = db.Put(randkv.GetTestKey(22), randkv.RandomValue(24))
	assert.Nil(t, err)
	err = db.Put(randkv.GetTestKey(22), randkv.RandomValue(24))
	val3, err := db.Get(randkv.GetTestKey(22))
	assert.Nil(t, err)
	assert.NotNil(t, val3)

	// 4.值被删除后再 Get
	err = db.Put(randkv.GetTestKey(33), randkv.RandomValue(24))
	assert.Nil(t, err)
	err = db.Delete(randkv.GetTestKey(33))
	assert.Nil(t, err)
	val4, err := db.Get(randkv.GetTestKey(33))
	assert.Equal(t, 0, len(val4))
	assert.Equal(t, flydb.ErrKeyNotFound, err)

	// 5.转换为了旧的数据文件，从旧的数据文件上获取 value
	for i := 100; i < 1000000; i++ {
		err := db.Put(randkv.GetTestKey(i), randkv.RandomValue(128))
		assert.Nil(t, err)
	}
	assert.Equal(t, 2, len(db.olderFiles))
	val5, err := db.Get(randkv.GetTestKey(101))
	assert.Nil(t, err)
	assert.NotNil(t, val5)

	// 6.重启后，前面写入的数据都能拿到
	err = db.Close()
	assert.Nil(t, err)

	// 重启数据库
	db2, err := NewDB(opts)
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
	assert.Equal(t, flydb.ErrKeyNotFound, err)
}

func TestDB_Delete(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-delete")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// 1.正常删除一个存在的 key
	err = db.Put(randkv.GetTestKey(11), randkv.RandomValue(128))
	assert.Nil(t, err)
	err = db.Delete(randkv.GetTestKey(11))
	assert.Nil(t, err)
	_, err = db.Get(randkv.GetTestKey(11))
	assert.Equal(t, flydb.ErrKeyNotFound, err)

	// 2.删除一个不存在的 key
	err = db.Delete([]byte("unknown key"))
	assert.Nil(t, err)

	// 3.删除一个空的 key
	err = db.Delete(nil)
	assert.Equal(t, flydb.ErrKeyIsEmpty, err)

	// 4.值被删除之后重新 Put
	err = db.Put(randkv.GetTestKey(22), randkv.RandomValue(128))
	assert.Nil(t, err)
	err = db.Delete(randkv.GetTestKey(22))
	assert.Nil(t, err)

	err = db.Put(randkv.GetTestKey(22), randkv.RandomValue(128))
	assert.Nil(t, err)
	val1, err := db.Get(randkv.GetTestKey(22))
	assert.NotNil(t, val1)
	assert.Nil(t, err)

	// 5.重启之后，再进行校验
	err = db.Close()
	assert.Nil(t, err)

	// 重启数据库
	db2, err := NewDB(opts)
	_, err = db2.Get(randkv.GetTestKey(11))
	assert.Equal(t, flydb.ErrKeyNotFound, err)

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
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// 数据库为空
	keys1 := db.GetListKeys()
	assert.Equal(t, 0, len(keys1))

	// 只有一条数据
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
	defer destroyDB(db)
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
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put(randkv.GetTestKey(10), randkv.GetTestKey(10))
	assert.Nil(t, err)

	err = db.Sync()
	assert.Nil(t, err)
}
