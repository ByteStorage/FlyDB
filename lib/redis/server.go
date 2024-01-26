package redis

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/structure"
	"log"
	"sync"

	"github.com/tidwall/redcon"
)

type FlyDBServer struct {
	Dbs    map[int]interface{}
	Server *redcon.Server
	lock   sync.RWMutex
}

// Listen starts listening for incoming connections.
func (svr *FlyDBServer) Listen() {
	log.Println("FlyDB-Redis Server Start Success On: ", config.DefaultRedisAddr)
	_ = svr.Server.ListenAndServe()
}

// Accept is called when a new client connects.
func (svr *FlyDBServer) Accept(conn redcon.Conn) bool {
	svr.lock.Lock()
	defer svr.lock.Unlock()

	cli := new(FlyDBClient)
	cli.Server = svr
	cli.DB = svr.Dbs
	conn.SetContext(cli)
	return true
}

// Close closes the server.
func (svr *FlyDBServer) Close(conn redcon.Conn, err error) {
	if db, ok := svr.Dbs[0].(*structure.StringStructure); ok {
		db.Clean()
	} else if dbh, ok := svr.Dbs[1].(*structure.HashStructure); ok {
		dbh.Clean()
	}
	_ = svr.Server.Close()
	log.Println("FlyDB-Redis Server Stop Success On: ", config.DefaultRedisAddr)
	return
}

// StartRedisServer starts a Redis server.
func StartRedisServer() {
	// open Redis data structure service
	options := config.DefaultOptions
	options.DirPath = config.RedisStringDirPath
	stringStructure, err := structure.NewStringStructure(options)
	if err != nil {
		panic(err)
	}

	options.DirPath = config.RedisHashDirPath
	hashStructure, err := structure.NewHashStructure(options)
	if err != nil {
		panic(err)
	}

	// initialize FlyDBServer
	flydbServer := FlyDBServer{
		Dbs: make(map[int]interface{}),
	}
	flydbServer.Dbs[0] = stringStructure
	flydbServer.Dbs[1] = hashStructure

	// initialize a Redis server
	flydbServer.Server = redcon.NewServer(config.DefaultRedisAddr,
		ClientCommands, flydbServer.Accept, flydbServer.Close)
	flydbServer.Listen()
}
