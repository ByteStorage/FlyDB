package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ByteStorage/FlyDB/lib/proto/ghash"
	"github.com/ByteStorage/FlyDB/lib/proto/glist"
	"github.com/ByteStorage/FlyDB/lib/proto/gset"
	"github.com/ByteStorage/FlyDB/lib/proto/gstring"
	"github.com/ByteStorage/FlyDB/lib/proto/gzset"
)

// Client is a grpc client
type Client struct {
	Addr                 string // db server address
	conn                 *grpc.ClientConn
	gStringServiceClient gstring.GStringServiceClient
	gHashServiceClient   ghash.GHashServiceClient
	gListServiceClient   glist.GListServiceClient
	gSetServiceClient    gset.GSetServiceClient
	gZSetServiceClient   gzset.GZSetServiceClient
}

func NewClient(addr string) (*Client, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)

	if conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return nil, err
	}

	return &Client{
		Addr: addr,
		conn: conn,
	}, nil
}

func (c *Client) getGrpcConn() (*grpc.ClientConn, error) {
	if c.conn == nil {
		conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}

		c.conn = conn
	}

	return c.conn, nil
}

// newGrpcClient returns a grpc client
func (c *Client) newGrpcClient() (gstring.GStringServiceClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)

	if c.gStringServiceClient != nil {
		return c.gStringServiceClient, nil
	}

	if conn, err = c.getGrpcConn(); err != nil {
		return nil, err
	}

	c.gStringServiceClient = gstring.NewGStringServiceClient(conn)

	return c.gStringServiceClient, nil
}

func (c *Client) newHashGrpcClient() (ghash.GHashServiceClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)

	if c.gHashServiceClient != nil {
		return c.gHashServiceClient, nil
	}

	if conn, err = c.getGrpcConn(); err != nil {
		return nil, err
	}

	c.gHashServiceClient = ghash.NewGHashServiceClient(conn)

	return c.gHashServiceClient, nil
}

func (c *Client) newListGrpcClient() (glist.GListServiceClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)

	if c.gListServiceClient != nil {
		return c.gListServiceClient, nil
	}

	if conn, err = c.getGrpcConn(); err != nil {
		return nil, err
	}

	c.gListServiceClient = glist.NewGListServiceClient(conn)

	return c.gListServiceClient, nil
}

func (c *Client) newSetGrpcClient() (gset.GSetServiceClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)

	if c.gSetServiceClient != nil {
		return c.gSetServiceClient, nil
	}

	if conn, err = c.getGrpcConn(); err != nil {
		return nil, err
	}

	c.gSetServiceClient = gset.NewGSetServiceClient(conn)

	return c.gSetServiceClient, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}
