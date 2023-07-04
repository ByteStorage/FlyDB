package cmd

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/ByteStorage/FlyDB/flydb"
	"github.com/desertbit/grumble"
	"os"
)

var db *engine.DB

func startServer(c *grumble.Context) error {
	if len(c.Args) == 0 {
		options := config.DefaultOptions
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
		if db != nil {
			err := db.Close()
			if err != nil {
				fmt.Println("flydb clean error: ", err)
				return err
			}
			db = nil
		}
		err := os.RemoveAll("FlyDB-cmd")
		if err != nil {
			fmt.Println("flydb clean error: ", err)
			return err
		}
		err = os.RemoveAll("/tmp/.FlyDB_Cli.history")
		fmt.Println("flydb clean success")
	}
	return nil
}

