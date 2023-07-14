package client

import (
	"context"
	"errors"
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
	client := gstring.NewGStringServiceClient(conn)
	return client, nil
}

// Put puts a key-value pair into the db by client api
func (c *Client) Put(key string, value interface{}) error {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	req := &gstring.SetRequest{Key: key}
	switch v := value.(type) {
	case string:
		req.Value = &gstring.SetRequest_StringValue{StringValue: v}
	case int32:
		req.Value = &gstring.SetRequest_Int32Value{Int32Value: v}
	case int64:
		req.Value = &gstring.SetRequest_Int64Value{Int64Value: v}
	case float32:
		req.Value = &gstring.SetRequest_Float32Value{Float32Value: v}
	case float64:
		req.Value = &gstring.SetRequest_Float64Value{Float64Value: v}
	case bool:
		req.Value = &gstring.SetRequest_BoolValue{BoolValue: v}
	case []byte:
		req.Value = &gstring.SetRequest_BytesValue{BytesValue: v}
	default:
		return errors.New("unknown value type")
	}
	put, err := client.Put(context.Background(), req)
	if err != nil {
		return errors.New("client put failed: " + err.Error())
	}
	if !put.Ok {
		return errors.New("put failed")
	}
	return nil
}

// Get gets a value by key from the db by client api
func (c *Client) Get(key string) (interface{}, error) {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return nil, err
	}
	get, err := client.Get(context.Background(), &gstring.GetRequest{Key: key})
	if err != nil {
		return nil, err
	}
	switch get.Value.(type) {
	case *gstring.GetResponse_StringValue:
		return get.Value.(*gstring.GetResponse_StringValue).StringValue, nil
	case *gstring.GetResponse_Int32Value:
		return get.Value.(*gstring.GetResponse_Int32Value).Int32Value, nil
	case *gstring.GetResponse_Int64Value:
		return get.Value.(*gstring.GetResponse_Int64Value).Int64Value, nil
	case *gstring.GetResponse_Float32Value:
		return get.Value.(*gstring.GetResponse_Float32Value).Float32Value, nil
	case *gstring.GetResponse_Float64Value:
		return get.Value.(*gstring.GetResponse_Float64Value).Float64Value, nil
	case *gstring.GetResponse_BoolValue:
		return get.Value.(*gstring.GetResponse_BoolValue).BoolValue, nil
	case *gstring.GetResponse_BytesValue:
		return get.Value.(*gstring.GetResponse_BytesValue).BytesValue, nil
	default:
		return nil, errors.New("get failed")
	}
}

// Del deletes a key-value pair from the db by client api
func (c *Client) Del(key string) error {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return err
	}
	del, err := client.Del(context.Background(), &gstring.DelRequest{Key: key})
	if err != nil {
		return err
	}
	if !del.Ok {
		return errors.New("del failed")
	}
	return nil
}

func (c *Client) Type(key string) (string, error) {
	panic("implement me")
}
