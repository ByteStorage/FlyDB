package cmd

import (
	"errors"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/ByteStorage/FlyDB/flydb"
	"github.com/desertbit/grumble"
	"os"
)

var db *engine.DB
var options = config.DefaultOptions

func startServer(c *grumble.Context) error {
	if len(c.Args) == 0 {
		options.DirPath = "FlyDB-cmd"
		err := os.Mkdir(options.DirPath, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			fmt.Println("flydb start error: ", err)
			return nil
		}

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
			fmt.Println("no server to stop")
			return nil
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
		err = os.RemoveAll(options.DirPath)
		if err != nil {
			fmt.Println("flydb clean error: ", err)
			return err
		}
		db = nil
		fmt.Println("flydb clean success")
	}
	return nil
}
