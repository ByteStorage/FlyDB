package engine

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func BenchmarkDB_Put(b *testing.B) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-benchmark")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer destroyDB(db)
	assert.Nil(b, err)
	assert.NotNil(b, db)

	for n := 0; n < b.N; n++ {
		err = db.Put(randkv.GetTestKey(n), randkv.RandomValue(24))
		assert.Nil(b, err)
	}
}

func TestPutAndGet(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-benchmark")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	start := time.Now()
	for n := 0; n < 500000; n++ {
		err = db.Put(randkv.GetTestKey(n), randkv.RandomValue(24))
		assert.Nil(t, err)
	}
	end := time.Now()
	fmt.Println("put time: ", end.Sub(start).String())

	start = time.Now()
	for n := 0; n < 500000; n++ {
		_, err = db.Get(randkv.GetTestKey(n))
		assert.Nil(t, err)
	}
	end = time.Now()
	fmt.Println("get time: ", end.Sub(start).String())
}
