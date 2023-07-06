package server

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine/grpc/service"
	"github.com/ByteStorage/FlyDB/flydb"
)

func StartServer() {
	options := config.DefaultOptions
	db, err := flydb.NewFlyDB(options)
	if err != nil {
		fmt.Println("flydb start error: ", err)
		return
	}
	s := service.NewService(config.DefaultAddr, db)
	s.StartServer()
}

func StopServer() {
	panic("implement me")
}

func CleanServer() {
	panic("implement me")
}
