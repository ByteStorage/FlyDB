package server

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine/grpc/service"
	"github.com/ByteStorage/FlyDB/structure"
	"os"
)

var s *service.Service

func StartServer() {
	options := config.DefaultOptions
	str, err := structure.NewStringStructure(options)
	if err != nil {
		fmt.Println("flydb start error: ", err)
		return
	}
	s = service.NewService(config.DefaultAddr, str)
	s.StartServer()
}

func CleanServer() {
	err := os.Remove(config.DefaultOptions.DirPath)
	if err != nil {
		fmt.Println("flydb clean error: ", err)
		return
	}
}
