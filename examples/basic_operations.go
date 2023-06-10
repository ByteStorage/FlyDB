package main

import (
	"fmt"
	"github.com/ByteStorage/flydb"
	"github.com/ByteStorage/flydb/config"
)

func main() {
	options := config.DefaultOptions
	options.DirPath = "/tmp/flydb"
	db, err := flydb.NewFlyDB(options)
	if err != nil {
		panic(err)
	}

	err = db.Put([]byte("name"), []byte("flydb-example"))
	if err != nil {
		panic(err)
	}

	val, err := db.Get([]byte("name"))
	if err != nil {
		panic(err)
	}

	fmt.Println("name value => ", string(val))

	err = db.Delete([]byte("name"))
	if err != nil {
		panic(err)
	}

}
