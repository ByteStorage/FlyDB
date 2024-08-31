package server

import (
	"fmt"
	"os"

	"github.com/ByteStorage/FlyDB/config"
	base "github.com/ByteStorage/FlyDB/db/grpc"
)

func StartServer() {
	options := config.DefaultOptions
	options.FIOType = config.MmapIOType
	service, err := base.NewService(options, "0.0.0.0:8999")
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
