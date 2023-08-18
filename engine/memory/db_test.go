package memory

import (
	_ "net/http/pprof"
)

//func TestPutAndGet(t *testing.T) {
//	opts := config.DefaultOptions
//	dir, _ := os.MkdirTemp("", "flydb-benchmark")
//	opts.DirPath = dir
//	opts.DataFileSize = 64 * 1024 * 1024
//	db, err := NewDbWal(opts)
//	defer db.Clean()
//	assert.Nil(t, err)
//	assert.NotNil(t, db)
//
//	start := time.Now()
//	for n := 0; n < 500000; n++ {
//		err = db.PutByWal(randkv.GetTestKey(n), randkv.RandomValue(24))
//		assert.Nil(t, err)
//	}
//	end := time.Now()
//	fmt.Println("put time: ", end.Sub(start).String())
//
//	start = time.Now()
//	for n := 0; n < 500000; n++ {
//		_, err = db.GetByWal(randkv.GetTestKey(n))
//		assert.Nil(t, err)
//	}
//	end = time.Now()
//	fmt.Println("get time: ", end.Sub(start).String())
//}
