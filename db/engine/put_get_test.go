package engine

import (
	"os"
	"testing"

	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
)

func BenchmarkPut(b *testing.B) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-benchmark")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(b, err)
	assert.NotNil(b, db)

	n := 0
	for b.Loop() {
		err = db.Put(randkv.GetTestKey(n), randkv.RandomValue(24))
		assert.Nil(b, err)
		n++
	}
}

func BenchmarkGet(b *testing.B) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-benchmark")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := NewDB(opts)
	defer db.Clean()
	assert.Nil(b, err)
	assert.NotNil(b, db)

	for n := 0; n < b.N; n++ {
		err = db.Put(randkv.GetTestKey(n), randkv.RandomValue(24))
		assert.Nil(b, err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, err = db.Get(randkv.GetTestKey(n))
		assert.Nil(b, err)
	}
}
