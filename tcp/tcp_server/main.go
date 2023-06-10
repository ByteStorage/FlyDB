package main

import (
	"fmt"
	"github.com/ByteStorage/flydb/config"
	"github.com/ByteStorage/flydb/tcp"
	"strconv"
)

func main() {
	tpcDefaultConfig := config.Init()
	err := tcp.ListenAndServeBySignal(&tcp.Config{
		Address: tpcDefaultConfig.Host + ":" + strconv.Itoa(tpcDefaultConfig.Port),
	}, tcp.NewHandler())
	if err != nil {
		fmt.Println(err)
	}
}
