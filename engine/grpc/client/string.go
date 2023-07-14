package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/ByteStorage/FlyDB/lib/proto/gstring"
)

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

func (c *Client) StrLen(key string) (int32, error) {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return 0, nil
	}
	resp, err := client.StrLen(context.Background(), &gstring.StrLenRequest{Key: key})
	if err != nil {
		return 0, err
	}
	return resp.Length, nil
}

func (c *Client) Append(key string, value string) error {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return err
	}
	resp, err := client.Append(context.Background(), &gstring.AppendRequest{Key: key, Value: value})
	if err != nil {
		return err
	}
	if !resp.Ok {
		return errors.New("Append operation failed")
	}
	return nil
}

func (c *Client) GetSet(key string, value interface{}) (interface{}, error) {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error()), nil
	}
	req := &gstring.GetSetRequest{
		Key: key,
	}

	switch v := value.(type) {
	case string:
		req.Value = &gstring.GetSetRequest_StringValue{StringValue: v}
	case int32:
		req.Value = &gstring.GetSetRequest_Int32Value{Int32Value: v}
	case int64:
		req.Value = &gstring.GetSetRequest_Int64Value{Int64Value: v}
	case float32:
		req.Value = &gstring.GetSetRequest_Float32Value{Float32Value: v}
	case float64:
		req.Value = &gstring.GetSetRequest_Float64Value{Float64Value: v}
	case bool:
		req.Value = &gstring.GetSetRequest_BoolValue{BoolValue: v}
	case []byte:
		req.Value = &gstring.GetSetRequest_BytesValue{BytesValue: v}
	default:
		return nil, fmt.Errorf("unknown value type")
	}

	resp, err := client.GetSet(context.Background(), req)
	if err != nil {
		return nil, err
	}

	switch v := resp.Value.(type) {
	case *gstring.GetSetResponse_StringValue:
		return v.StringValue, nil
	case *gstring.GetSetResponse_Int32Value:
		return v.Int32Value, nil
	case *gstring.GetSetResponse_Int64Value:
		return v.Int64Value, nil
	case *gstring.GetSetResponse_Float32Value:
		return v.Float32Value, nil
	case *gstring.GetSetResponse_Float64Value:
		return v.Float64Value, nil
	case *gstring.GetSetResponse_BoolValue:
		return v.BoolValue, nil
	case *gstring.GetSetResponse_BytesValue:
		return v.BytesValue, nil
	default:
		return nil, fmt.Errorf("unknown response value type")
	}
}

func (c *Client) Exists(key string) (bool, error) {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return false, nil
	}
	resp, err := client.Exists(context.Background(), &gstring.ExistsRequest{Key: key})
	if err != nil {
		return false, nil
	}
	return resp.Exists, nil
}

func (c *Client) Expire(key string, ttl int64) error {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return err
	}
	expire, err := client.Expire(context.Background(), &gstring.ExpireRequest{Key: key, Expire: ttl})
	if err != nil {
		return err
	}
	if !expire.Ok {
		return errors.New("outof date")
	}
	return nil
}

func (c *Client) Persist(key string) error {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return err
	}
	persist, err := client.Persist(context.Background(), &gstring.PersistRequest{Key: key})
	if err != nil {
		return err
	}
	if !persist.Ok {
		return errors.New("Not persist")
	}
	return nil
}

func (c *Client) MGet(keys []string) ([]interface{}, error) {
	client, err := newGrpcClient(c.Addr)
	if err != nil {
		return nil, err
	}

	resp, err := client.MGet(context.Background(), &gstring.MGetRequest{Keys: keys})
	if err != nil {
		return nil, err
	}
	var values []interface{}
	for _, value := range resp.Values {
		switch v := value.Value.(type) {
		case *gstring.MGetValue_StringValue:
			values = append(values, v.StringValue)
		case *gstring.MGetValue_Int32Value:
			values = append(values, v.Int32Value)
		case *gstring.MGetValue_Int64Value:
			values = append(values, v.Int64Value)
		case *gstring.MGetValue_Float32Value:
			values = append(values, v.Float32Value)
		case *gstring.MGetValue_Float64Value:
			values = append(values, v.Float64Value)
		case *gstring.MGetValue_BoolValue:
			values = append(values, v.BoolValue)
		case *gstring.MGetValue_BytesValue:
			values = append(values, v.BytesValue)
		}
	}
	return values, nil
}
