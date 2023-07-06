package client

import (
	"context"
	"errors"
	"github.com/ByteStorage/FlyDB/lib/proto/dbs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is a grpc client
type Client struct {
	Addr string // db server address
}

// newGrpcClient returns a grpc client
func newGrpcClient(addr string) (dbs.FlyDBServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := dbs.NewFlyDBServiceClient(conn)
	return client, nil
}

// Put puts a key-value pair into the db by client api
func (c *Client) Put(key []byte, value []byte) error {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return err
	}
	put, err := client.Put(context.Background(), &dbs.PutRequest{Key: key, Value: value})
	if err != nil {
		return err
	}
	if !put.Ok {
		return errors.New("put failed")
	}
	return nil
}

// Get gets a value by key from the db by client api
func (c *Client) Get(key []byte) ([]byte, error) {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return nil, err
	}
	get, err := client.Get(context.Background(), &dbs.GetRequest{Key: key})
	if err != nil {
		return nil, err
	}
	return get.Value, nil
}

// Del deletes a key-value pair from the db by client api
func (c *Client) Del(key []byte) error {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return err
	}
	del, err := client.Del(context.Background(), &dbs.DelRequest{Key: key})
	if err != nil {
		return err
	}
	if !del.Ok {
		return errors.New("del failed")
	}
	return nil
}

// Keys gets all keys from the db by client api
func (c *Client) Keys() ([][]byte, error) {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return nil, err
	}
	keys, err := client.Keys(context.Background(), &dbs.KeysRequest{})
	if err != nil {
		return nil, err
	}
	result := make([][]byte, len(keys.Keys))
	for i, key := range keys.Keys {
		result[i] = key
	}
	return result, nil
}
