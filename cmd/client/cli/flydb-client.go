package main

import (
	"github.com/ByteStorage/FlyDB/cmd/client"
	"github.com/desertbit/grumble"
)

func main() {
	// start client CLI
	grumble.Main(client.App)
}
