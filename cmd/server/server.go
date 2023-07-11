package server

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine/grpc/service"
	"github.com/ByteStorage/FlyDB/structure"
)

var str *structure.StringStructure

func StartServer() {
	var err error
	options := config.DefaultOptions
	str, err = structure.NewStringStructure(options)
	if err != nil {
		fmt.Println("flydb start error: ", err)
		return
	}
	s := service.NewService(config.DefaultAddr, str)
	s.StartServer()
}

func StopServer() {
	if str == nil {
		fmt.Println("flydb stop error: ", "flydb not running")
		return
	}
	err := str.Stop()
	if err != nil {
		fmt.Println("flydb stop error: ", err)
		return
	}
}

func CleanServer() {
	if str == nil {
		fmt.Println("flydb clean error: ", "flydb not running")
		return
	}
	str.Clean()
}
