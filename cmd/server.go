package cmd

import (
	"errors"
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
	if len(c.Args) == 0 {
		if db == nil {
			err := errors.New("no server to stop")
			fmt.Println("flydb stop error: ", err)
			return err
		}
		err := db.Close()
		if err != nil {
			fmt.Println("flydb stop error: ", err)
			return err
		}
		fmt.Println("flydb stop success")
		db = nil
	}
	return nil
}

func cleanServer(c *grumble.Context) error {
	if len(c.Args) == 0 {
		if db == nil {
			err := errors.New("no server to clean")
			fmt.Println("flydb clean error: ", err)
			return err
		}
		err := db.Close()
		if err != nil {
			fmt.Println("flydb clean error: ", err)
			return err
		}
		err = db.Clean()
		if err != nil {
			fmt.Println("flydb clean error: ", err)
			return err
		}

		fmt.Println("flydb clean success")
	}
	return nil
}
