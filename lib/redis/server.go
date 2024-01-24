package redis

import (
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

func (svr *FlyDBServer) Listen() {
	log.Println("FlyDB-Redis Server running, ready to accept connections...")
	_ = svr.Server.ListenAndServe()
}

func (svr *FlyDBServer) Accept(conn redcon.Conn) bool {
	svr.lock.Lock()
	defer svr.lock.Unlock()

	cli := new(FlyDBClient)
	cli.Server = svr
	cli.DB = svr.Dbs
	conn.SetContext(cli)
	return true
}

func (svr *FlyDBServer) Close() error {
	if db, ok := svr.Dbs[0].(*structure.StringStructure); ok {
		db.Clean()
	} else if dbh, ok := svr.Dbs[1].(*structure.HashStructure); ok {
		dbh.Clean()
	}
	return svr.Server.Close()
}
