package cluster

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Message struct {
}

type Grpc interface {
	SendMessage(addr string, msg Message) error
}

type Client struct {
}

func (c *Client) SendMessageToSlave(addr string, msg Message) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	//如何用范型来解决这个问题？
	//if msg is proto.SlaveGetRequest
	//proto.NewSlaveClient(conn).Get()
	panic("implement me")
}
