package client

import (
	"github.com/ByteStorage/FlyDB/lib/proto/ghash"
	"github.com/ByteStorage/FlyDB/lib/proto/gstring"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is a grpc client
type Client struct {
	Addr string // db server address
}

// newGrpcClient returns a grpc client
func newGrpcClient(addr string) (gstring.GStringServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	//client := gstring.NewGStringServiceClient(conn)
	client := gstring.NewGStringServiceClient(conn)
	return client, nil
}

func newHashGrpcClient(addr string) (ghash.GHashServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	//client := gstring.NewGStringServiceClient(conn)
	client := ghash.NewGHashServiceClient(conn)
	return client, nil
}
