package main

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/redis"
	flydb_stru "github.com/ByteStorage/FlyDB/structure"
	"github.com/tidwall/redcon"
)

func main() {
	// 打开 Redis 数据结构服务
	stringStructure, err := flydb_stru.NewStringStructure(config.DefaultOptions)
	if err != nil {
		panic(err)
	}
	hashStructure, err := flydb_stru.NewHashStructure(config.DefaultOptions)
	if err != nil {
		panic(err)
	}

	// 初始化 FlyDBServer
	flydbServer := &redis.FlyDBServer{
		Dbs: make(map[int]interface{}),
	}
	flydbServer.Dbs[0] = stringStructure
	flydbServer.Dbs[1] = hashStructure

	// 初始化一个 Redis 服务端
	flydbServer.Server = redcon.NewServer(config.DefaultAddr,
		redis.ClientCommands, flydbServer.Accept, nil)
	flydbServer.Listen()
}
