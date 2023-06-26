package cmd

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/ByteStorage/FlyDB/flydb"
	"github.com/desertbit/grumble"
)

var db *engine.DB

func startServer(c *grumble.Context) error {
	if len(c.Args) == 0 {
		options := config.DefaultOptions
		flyDb, err := flydb.NewFlyDB(options)
		if err != nil {
			fmt.Println("flydb start error: ", err)
			return err
		}
		db = flyDb
		fmt.Println("flydb start success")
	}
	return nil
}

func stopServer(c *grumble.Context) error {
	panic("implement me")
}

func cleanServer(c *grumble.Context) error {
	panic("implement me")
}
