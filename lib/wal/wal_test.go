package wal

import (
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWal_Put(t *testing.T) {
	opt := Options{
		DirPath:  "./wal_test",
		LogNum:   100,
		FileSize: 100 * 1024 * 1024,
		SaveTime: 100 * 1000,
	}
	wal, err := NewWal(opt)
	defer wal.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, wal)
	start := time.Now()
	for n := 0; n < 500000; n++ {
		err = wal.Put(randkv.GetTestKey(n), randkv.RandomValue(24))
		assert.Nil(t, err)
	}
	end := time.Now()
	t.Log("put time: ", end.Sub(start).String())
}
