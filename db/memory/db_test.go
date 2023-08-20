package memory

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestPutAndGet(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-benchmark")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	memoryOptions := config.DbMemoryOptions{
		Option:       opts,
		LogNum:       100,
		FileSize:     100 * 1024 * 1024,
		SaveTime:     100 * 1000,
		MemSize:      256 * 1024 * 1024,
		TotalMemSize: 2 * 1024 * 1024 * 1024,
	}
	db, err := NewDB(memoryOptions)
	defer db.Clean()
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
