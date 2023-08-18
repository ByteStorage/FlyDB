package compress

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
	db, err := NewDbCompress(opts)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	start := time.Now()
	err = db.Put([]byte("key"), randkv.RandomValue(10000000))
	assert.Nil(t, err)

	end := time.Now()
	fmt.Println("put time: ", end.Sub(start).String())

	start = time.Now()
	_, err = db.Get([]byte("key"))
	assert.Nil(t, err)
	end = time.Now()
	fmt.Println("get time: ", end.Sub(start).String())
}
