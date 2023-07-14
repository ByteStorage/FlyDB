package client

import (
	"context"
	"errors"
	"github.com/ByteStorage/FlyDB/lib/proto/ghash"
)

// HSet puts a key-value pair into the db by client api
func (c *Client) HSet(key, field string, value interface{}) error {
	client, err := newHashGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	req := &ghash.GHashSetRequest{Key: key, Field: field}
	err = setValue(req, value)
	if err != nil {
		return err
	}

	put, err := client.HSet(context.Background(), req)
	if err != nil {
		return errors.New("client put failed: " + err.Error())
	}
	if !put.Ok {
		return errors.New("put failed")
	}
	return nil
}

// HGet gets a value by key from the db by client api
func (c *Client) HGet(key, field string) (interface{}, error) {
	client, err := newHashGrpcClient(c.Addr)
	if err != nil {
		return nil, err
	}
	get, err := client.HGet(context.Background(),
		&ghash.GHashGetRequest{Key: key, Field: field})
	if err != nil {
		return nil, err
	}
	value, err := getValue(get)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// HDel deletes a key-value pair from the db by client api
func (c *Client) HDel(key, field string) error {
	client, err := newHashGrpcClient(c.Addr)
	if err != nil {
		return err
	}
	del, err := client.HDel(context.Background(),
		&ghash.GHashDelRequest{Key: key, Field: field})
	if err != nil {
		return err
	}
	if !del.Ok {
		return errors.New("del failed")
	}

	return nil
}

func setValue(req *ghash.GHashSetRequest, value interface{}) error {
	switch v := value.(type) {
	case string:
		req.Value = &ghash.GHashSetRequest_StringValue{StringValue: v}
	case int32:
		req.Value = &ghash.GHashSetRequest_Int32Value{Int32Value: v}
	case int64:
		req.Value = &ghash.GHashSetRequest_Int64Value{Int64Value: v}
	case float32:
		req.Value = &ghash.GHashSetRequest_Float32Value{Float32Value: v}
	case float64:
		req.Value = &ghash.GHashSetRequest_Float64Value{Float64Value: v}
	case bool:
		req.Value = &ghash.GHashSetRequest_BoolValue{BoolValue: v}
	case []byte:
		req.Value = &ghash.GHashSetRequest_BytesValue{BytesValue: v}
	default:
		return errors.New("unknown value type")
	}
	return nil
}

func getValue(resp *ghash.GHashGetResponse) (interface{}, error) {
	switch resp.Value.(type) {
	case *ghash.GHashGetResponse_StringValue:
		return resp.Value.(*ghash.GHashGetResponse_StringValue).StringValue, nil
	case *ghash.GHashGetResponse_Int32Value:
		return resp.Value.(*ghash.GHashGetResponse_Int32Value).Int32Value, nil
	case *ghash.GHashGetResponse_Int64Value:
		return resp.Value.(*ghash.GHashGetResponse_Int64Value).Int64Value, nil
	case *ghash.GHashGetResponse_Float32Value:
		return resp.Value.(*ghash.GHashGetResponse_Float32Value).Float32Value, nil
	case *ghash.GHashGetResponse_Float64Value:
		return resp.Value.(*ghash.GHashGetResponse_Float64Value).Float64Value, nil
	case *ghash.GHashGetResponse_BoolValue:
		return resp.Value.(*ghash.GHashGetResponse_BoolValue).BoolValue, nil
	case *ghash.GHashGetResponse_BytesValue:
		return resp.Value.(*ghash.GHashGetResponse_BytesValue).BytesValue, nil
	default:
		return nil, errors.New("unknown value type")
	}
}
