package main

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	tcp2 "github.com/ByteStorage/FlyDB/protocol/tcp"
	"strconv"
)

func main() {
	tpcDefaultConfig := config.Init()
	err := tcp2.ListenAndServeBySignal(&tcp2.Config{
		Address: tpcDefaultConfig.Host + ":" + strconv.Itoa(tpcDefaultConfig.Port),
	}, tcp2.NewHandler())
	if err != nil {
		fmt.Println(err)
	}
}
