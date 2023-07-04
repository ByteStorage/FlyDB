package server

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/flydb"
)

func StartServer() {
	options := config.DefaultOptions
	_, err := flydb.NewFlyDB(options)
	if err != nil {
		fmt.Println("flydb start error: ", err)
		return
	}
	fmt.Println("flydb start success")
}

func StopServer() {
	panic("implement me")
}

func CleanServer() {
	panic("implement me")
}
