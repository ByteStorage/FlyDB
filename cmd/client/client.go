package client

import (
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
		panic(err)
	}

	return cliClient
}

func cliClientClose() error {
	if cliClient == nil {
		return nil
	}

	if err := cliClient.Close(); err != nil {
		return err
	}

	cliClient = nil

	return nil
}
