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
	err := os.Mkdir("./flydb-benchmark", os.ModePerm)
	memOpt := config.DefaultDbMemoryOptions
	memOpt.LogNum = 100
	memOpt.FileSize = 256 * 1024 * 1024
	memOpt.TotalMemSize = 2 * 1024 * 1024 * 1024
	memOpt.Option.DirPath = "./flydb-benchmark"

	db, err := NewDB(memOpt)
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

func TestDb_Keys(t *testing.T) {
	err := os.Mkdir("./flydb", os.ModePerm)
	memOpt := config.DefaultDbMemoryOptions
	memOpt.LogNum = 100
	memOpt.FileSize = 256 * 1024 * 1024
	memOpt.TotalMemSize = 2 * 1024 * 1024 * 1024
	memOpt.Option.DirPath = "./"

	db, err := NewDB(memOpt)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	for n := 0; n < 100; n++ {
		err = db.Put(randkv.GetTestKey(n), randkv.RandomValue(24))
		assert.Nil(t, err)
	}

	keys, err := db.Keys()
	assert.Nil(t, err)
	assert.Equal(t, 100, len(keys))
	t.Log(keys)
}
