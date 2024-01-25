package main

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/redis"
	flydb_stru "github.com/ByteStorage/FlyDB/structure"
	"github.com/tidwall/redcon"
)

func main() {
	// open Redis data structure service
	options := config.DefaultOptions
	options.DirPath = config.RedisStringDirPath
	stringStructure, err := flydb_stru.NewStringStructure(options)
	if err != nil {
		panic(err)
	}

	options.DirPath = config.RedisHashDirPath
	hashStructure, err := flydb_stru.NewHashStructure(options)
	if err != nil {
		panic(err)
	}

	// initialize FlyDBServer
	flydbServer := &redis.FlyDBServer{
		Dbs: make(map[int]interface{}),
	}
	flydbServer.Dbs[0] = stringStructure
	flydbServer.Dbs[1] = hashStructure

	// initialize a Redis server
	flydbServer.Server = redcon.NewServer(config.DefaultAddr,
		redis.ClientCommands, flydbServer.Accept, nil)
	flydbServer.Listen()
}
