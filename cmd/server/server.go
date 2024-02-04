package server

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	base "github.com/ByteStorage/FlyDB/db/grpc"
	"github.com/ByteStorage/FlyDB/lib/redis"
	"os"
	"sync"
)

func StartServer() {
	var wg sync.WaitGroup

	// start flydb server
	wg.Add(1)
	go func() {
		defer wg.Done()
		options := config.DefaultOptions
		options.FIOType = config.MmapIOType
		service, err := base.NewService(options, config.DefaultAddr)
		if err != nil {
			fmt.Println("flydb start error: ", err)
			return
		}
		service.StartGrpcServer()
	}()

	// start flydb-redis server
	wg.Add(1)
	go func() {
		defer wg.Done()
		redis.StartRedisServer()
	}()

	// wait for signal
	wg.Wait()
}

func CleanServer() {
	err := os.Remove(config.DefaultOptions.DirPath)
	if err != nil {
		fmt.Println("flydb clean error: ", err)
		return
	}
}
