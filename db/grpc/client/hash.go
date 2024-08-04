package client

import (
	"context"
	"errors"

	"github.com/ByteStorage/FlyDB/lib/proto/ghash"
)

// HSet puts a key-value pair into the db by client api
func (c *Client) HSet(key, field string, value interface{}) error {
	client, err := c.newHashGrpcClient()
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
	client, err := c.newHashGrpcClient()
	if err != nil {
		return nil, err
	}
	get, err := client.HGet(context.Background(),
		&ghash.GHashGetRequest{Key: key, Field: field})
	if err != nil {
		return nil, err
	}

	if get.Value == nil {
		return nil, nil
	}

	value, err := getValue(get)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// HDel deletes a key-value pair from the db by client api
func (c *Client) HDel(key, field string) error {
	client, err := c.newHashGrpcClient()
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

func (c *Client) HExists(key, field string) (bool, error) {
	client, err := c.newHashGrpcClient()
	if err != nil {
		return false, err
	}

	exists, err := client.HExists(context.Background(),
		&ghash.GHashExistsRequest{Key: key, Field: field})
	if err != nil {
		return false, err
	}

	return exists.Ok, nil
}

func (c *Client) HLen(key string) (int64, error) {
	client, err := c.newHashGrpcClient()
	if err != nil {
		return 0, err
	}

	lenResp, err := client.HLen(context.Background(),
		&ghash.GHashLenRequest{Key: key})
	if err != nil {
		return 0, err
	}

	return lenResp.Length, nil
}

func (c *Client) HUpdate(key, field string, value interface{}) error {
	client, err := c.newHashGrpcClient()
	if err != nil {
		return err
	}
	req := &ghash.GHashUpdateRequest{Key: key, Field: field}
	err = updateValue(req, value)
	if err != nil {
		return err
	}

	update, err := client.HUpdate(context.Background(), req)
	if err != nil {
		return errors.New("client HSetNX failed: " + err.Error())
	}
	if !update.Ok {
		return errors.New("HSet failed")
	}

	return nil
}

func updateValue(req *ghash.GHashUpdateRequest, value interface{}) error {
	switch v := value.(type) {
	case string:
		req.Value = &ghash.GHashUpdateRequest_StringValue{StringValue: v}
	case int32:
		req.Value = &ghash.GHashUpdateRequest_Int32Value{Int32Value: v}
	case int64:
		req.Value = &ghash.GHashUpdateRequest_Int64Value{Int64Value: v}
	case float32:
		req.Value = &ghash.GHashUpdateRequest_Float32Value{Float32Value: v}
	case float64:
		req.Value = &ghash.GHashUpdateRequest_Float64Value{Float64Value: v}
	case bool:
		req.Value = &ghash.GHashUpdateRequest_BoolValue{BoolValue: v}
	case []byte:
		req.Value = &ghash.GHashUpdateRequest_BytesValue{BytesValue: v}
	default:
		return errors.New("unknown value type")
	}
	return nil
}

func (c *Client) HIncrBy(key, field string, value int64) (int64, error) {
	client, err := c.newHashGrpcClient()
	if err != nil {
		return 0, err
	}

	incrResp, err := client.HIncrBy(context.Background(),
		&ghash.GHashIncrByRequest{Key: key, Field: field, Value: value})
	if err != nil {
		return 0, err
	}

	return incrResp.Value, nil
}

func (c *Client) HIncrByFloat(key, field string, value float64) (float64, error) {
	client, err := c.newHashGrpcClient()
	if err != nil {
		return 0.0, err
	}

	incrResp, err := client.HIncrByFloat(context.Background(),
		&ghash.GHashIncrByFloatRequest{Key: key, Field: field, Value: value})
	if err != nil {
		return 0.0, err
	}

	return incrResp.Value, nil
}

func (c *Client) HDecrBy(key, field string, value int64) (int64, error) {
	client, err := c.newHashGrpcClient()
	if err != nil {
		return 0, err
	}

	decrResp, err := client.HDecrBy(context.Background(),
		&ghash.GHashDecrByRequest{Key: key, Field: field, Value: value})
	if err != nil {
		return 0, err
	}

	return decrResp.Value, nil
}

func (c *Client) HStrLen(key, field string) (int64, error) {
	client, err := c.newHashGrpcClient()
	if err != nil {
		return 0, err
	}

	strLenResp, err := client.HStrLen(context.Background(),
		&ghash.GHashStrLenRequest{Key: key, Field: field})
	if err != nil {
		return 0, err
	}

	return strLenResp.Length, nil
}

func (c *Client) HMove(key, dest, field string) error {
	client, err := c.newHashGrpcClient()
	if err != nil {
		return err
	}

	_, err = client.HMove(context.Background(),
		&ghash.GHashMoveRequest{Key: key, Dest: dest, Field: field})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) HSetNX(key, field string, value interface{}) error {
	client, err := c.newHashGrpcClient()
	if err != nil {
		return err
	}
	req := &ghash.GHashSetNXRequest{Key: key, Field: field}
	err = setNXValue(req, value)
	if err != nil {
		return err
	}

	put, err := client.HSetNX(context.Background(), req)
	if err != nil {
		return errors.New("client HSetNX failed: " + err.Error())
	}
	if !put.Ok {
		return errors.New("HSet failed")
	}
	return nil
}

func setNXValue(req *ghash.GHashSetNXRequest, value interface{}) error {
	switch v := value.(type) {
	case string:
		req.Value = &ghash.GHashSetNXRequest_StringValue{StringValue: v}
	case int32:
		req.Value = &ghash.GHashSetNXRequest_Int32Value{Int32Value: v}
	case int64:
		req.Value = &ghash.GHashSetNXRequest_Int64Value{Int64Value: v}
	case float32:
		req.Value = &ghash.GHashSetNXRequest_Float32Value{Float32Value: v}
	case float64:
		req.Value = &ghash.GHashSetNXRequest_Float64Value{Float64Value: v}
	case bool:
		req.Value = &ghash.GHashSetNXRequest_BoolValue{BoolValue: v}
	case []byte:
		req.Value = &ghash.GHashSetNXRequest_BytesValue{BytesValue: v}
	default:
		return errors.New("unknown value type")
	}
	return nil
}

func (c *Client) HType(key, field string) (string, error) {
	client, err := c.newHashGrpcClient()
	if err != nil {
		return "", err
	}

	typeResp, err := client.HType(context.Background(), &ghash.GHashTypeRequest{Key: key, Field: field})
	if err != nil {
		return "", err
	}

	return typeResp.Type, nil
}
