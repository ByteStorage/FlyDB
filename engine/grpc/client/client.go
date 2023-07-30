package client

import (
	"github.com/ByteStorage/FlyDB/lib/proto/ghash"
	"github.com/ByteStorage/FlyDB/lib/proto/glist"
	"github.com/ByteStorage/FlyDB/lib/proto/gstring"
	"github.com/ByteStorage/FlyDB/lib/proto/gzset"
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
	client := gstring.NewGStringServiceClient(conn)
	return client, nil
}

func newHashGrpcClient(addr string) (ghash.GHashServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := ghash.NewGHashServiceClient(conn)
	return client, nil
}

func newListGrpcClient(addr string) (glist.GListServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := glist.NewGListServiceClient(conn)
	return client, nil
}

func newZSetGrpcClient(addr string) (gzset.GZSetServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := gzset.NewGZSetServiceClient(conn)
	return client, nil
}
