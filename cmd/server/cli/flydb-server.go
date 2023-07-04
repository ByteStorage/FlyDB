package main

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/cmd/server"
	"os"
)

func main() {
	args := os.Args
	if len(args) >= 2 {
		fmt.Println("Usage: flydb-server [start|clean|stop]")
		return
	}
	if len(args) == 1 {
		//start server
		server.StartServer()
	}
	switch args[1] {
	case "start":
		server.StartServer()
	case "stop":
		server.StopServer()
	case "clean":
		server.CleanServer()
	}
}
