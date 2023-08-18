package engine

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestDB_WriteBatch(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-batch-1")
	opts.DirPath = dir
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// Do not submit data after writing
	wb := db.NewWriteBatch(config.DefaultWriteBatchOptions)
	err = wb.Put(randkv.GetTestKey(1), randkv.RandomValue(10))
	assert.Nil(t, err)
	err = wb.Delete(randkv.GetTestKey(2))
	assert.Nil(t, err)

	_, err = db.Get(randkv.GetTestKey(1))
	assert.Equal(t, _const.ErrKeyNotFound, err)

	// Normal submission data
	err = wb.Commit()
	assert.Nil(t, err)

	val, err := db.Get(randkv.GetTestKey(1))
	assert.NotNil(t, val)
	assert.Nil(t, err)

	wb2 := db.NewWriteBatch(config.DefaultWriteBatchOptions)
	err = wb2.Delete(randkv.GetTestKey(1))
	assert.Nil(t, err)
	err = wb2.Commit()
	assert.Nil(t, err)

	_, err = db.Get(randkv.GetTestKey(1))
	assert.Equal(t, _const.ErrKeyNotFound, err)
}

func TestDB_WriteBatchRestart(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-batch-2")
	opts.DirPath = dir
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put(randkv.GetTestKey(1), randkv.RandomValue(10))
	assert.Nil(t, err)

	wb := db.NewWriteBatch(config.DefaultWriteBatchOptions)
	err = wb.Put(randkv.GetTestKey(2), randkv.RandomValue(10))
	assert.Nil(t, err)
	err = wb.Delete(randkv.GetTestKey(1))
	assert.Nil(t, err)

	err = wb.Commit()
	assert.Nil(t, err)

	err = wb.Put(randkv.GetTestKey(3), randkv.RandomValue(10))
	assert.Nil(t, err)
	err = wb.Commit()
	assert.Nil(t, err)

	// Restart
	err = db.Close()
	assert.Nil(t, err)

	db2, err := NewDB(opts)
	time.Sleep(time.Millisecond * 100)
	defer db2.Clean()
	assert.Nil(t, err)

	_, err = db2.Get(randkv.GetTestKey(1))
	assert.Equal(t, _const.ErrKeyNotFound, err)

	// Judgment transaction sequence number
	assert.Equal(t, uint64(2), db.transSeqNo)
}

func TestDB_WriteBatch1(t *testing.T) {
	opts := config.DefaultOptions
	dir := "/tmp/batch-3"
	opts.DirPath = dir
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// Manual stop during batch submission
	wbopt := config.DefaultWriteBatchOptions
	wbopt.MaxBatchNum = 1000000
	wb := db.NewWriteBatch(wbopt)
	for i := 0; i < 500000; i++ {
		err = wb.Put(randkv.GetTestKey(i), randkv.RandomValue(1024))
		assert.Nil(t, err)
	}

	err = wb.Commit()
	assert.Nil(t, err)

}
