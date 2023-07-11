package main

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/cmd/server"
	"os"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage: flydb-server [start|clean]")
		return
	}
	if len(args) == 1 {
		//start server
		server.StartServer()
		return
	}
	switch args[1] {
	case "start":
		server.StartServer()
	case "clean":
		server.CleanServer()
	}
}
