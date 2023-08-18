package server

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	base "github.com/ByteStorage/FlyDB/db/grpc"
	"os"
)

func StartServer() {
	options := config.DefaultOptions
	options.FIOType = config.MmapIOType
	service, err := base.NewService(options, config.DefaultAddr)
	if err != nil {
		fmt.Println("flydb start error: ", err)
		return
	}
	service.StartGrpcServer()
}

func CleanServer() {
	err := os.Remove(config.DefaultOptions.DirPath)
	if err != nil {
		fmt.Println("flydb clean error: ", err)
		return
	}
}
