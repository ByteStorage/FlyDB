package client

import (
	"fmt"

	"github.com/ByteStorage/FlyDB/db/grpc/client"
)

var (
	Addr      string
	cliClient *client.Client
)

func newClient() *client.Client {
	var err error
	if cliClient != nil {
		return cliClient
	}

	if cliClient, err = client.NewClient(Addr); err != nil {
		fmt.Println("new client error: ", err)
	}

	return cliClient
}

func close() error {
	if cliClient == nil {
		return nil
	}

	if err := cliClient.Close(); err != nil {
		return err
	}

	cliClient = nil

	return nil
}
