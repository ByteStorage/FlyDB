package client

type Message struct {
}

type Grpc interface {
	SendMessage(addr string, msg Message) error
}

type Client struct {
}

func (c *Client) SendMessage(addr string, msg Message) error {
	panic("implement me")
}
