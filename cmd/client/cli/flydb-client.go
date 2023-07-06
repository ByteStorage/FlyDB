package main

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/cmd/client"
	"github.com/desertbit/grumble"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: flydb-cli [addr]")
		return
	}
	client.Addr = os.Args[1]
	os.Args = os.Args[:1]
	// start client CLI
	grumble.Main(client.App)
}
