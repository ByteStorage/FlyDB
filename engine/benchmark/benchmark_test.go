package benchmark

import (
	"bytes"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/ByteStorage/FlyDB/flydb"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"math/rand"
	"path/filepath"
	"testing"
	"time"
)

var FlyDB *engine.DB
var err error

const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().Unix())
}

func GetKey(n int) []byte {
	return []byte("test_key_" + fmt.Sprintf("%09d", n))
}

func GetValue() []byte {
	var str bytes.Buffer
	for i := 0; i < 512; i++ {
		str.WriteByte(alphabet[rand.Int()%36])
	}
	return []byte(str.String())
}

func init() {
	opts := config.DefaultOptions
	opts.DirPath = filepath.Join("benchmark", "flydbtest")

	FlyDB, err = flydb.NewFlyDB(opts)
	defer FlyDB.Clean()
	if err != nil {
		panic(err)
	}
}

func Benchmark_PutValue_FlyDB(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		err = FlyDB.Put(GetKey(n), GetValue())
		if err != nil {
			panic(err)
		}
	}
}

func Benchmark_GetValue_FlyDB(b *testing.B) {
	for i := 0; i < 500000; i++ {
		err = FlyDB.Put(GetKey(i), GetValue())
		if err != nil {
			panic(err)
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, err = FlyDB.Get(GetKey(n))
		if err != nil && err != _const.ErrKeyNotFound {
			panic(err)
		}
	}

}
