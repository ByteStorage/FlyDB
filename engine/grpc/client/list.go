package client

import (
	"context"
	"errors"
	"github.com/ByteStorage/FlyDB/lib/proto/glist"
)

func (c *Client) LPush(key string, value interface{}) error {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListLPushRequest{Key: key}

	switch v := value.(type) {
	case string:
		req.Value = &glist.GListLPushRequest_StringValue{StringValue: v}
	case int32:
		req.Value = &glist.GListLPushRequest_Int32Value{Int32Value: v}
	case int64:
		req.Value = &glist.GListLPushRequest_Int64Value{Int64Value: v}
	case float32:
		req.Value = &glist.GListLPushRequest_Float32Value{Float32Value: v}
	case float64:
		req.Value = &glist.GListLPushRequest_Float64Value{Float64Value: v}
	case bool:
		req.Value = &glist.GListLPushRequest_BoolValue{BoolValue: v}
	case []byte:
		req.Value = &glist.GListLPushRequest_BytesValue{BytesValue: v}
	default:
		return errors.New("unknown value type")
	}
	lpush, err := client.LPush(context.Background(), req)
	if err != nil {
		return errors.New("client LPush failed: " + err.Error())
	}
	if !lpush.Ok {
		return errors.New("LPush failed")
	}
	return nil
}

func (c *Client) LPushs(key string, values []interface{}) error {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListLPushsRequest{Key: key}

	// Convert the list of values to gRPC-compatible format
	for _, value := range values {
		switch v := value.(type) {
		case string:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_StringValue{StringValue: v}})
		case int32:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_Int32Value{Int32Value: v}})
		case int64:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_Int64Value{Int64Value: v}})
		case float32:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_Float32Value{Float32Value: v}})
		case float64:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_Float64Value{Float64Value: v}})
		case bool:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_BoolValue{BoolValue: v}})
		case []byte:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_BytesValue{BytesValue: v}})
		default:
			return errors.New("unknown value type")
		}
	}

	_, err = client.LPushs(context.Background(), req)
	if err != nil {
		return errors.New("client LPushs failed: " + err.Error())
	}

	return nil
}

func (c *Client) RPush(key string, value interface{}) error {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListRPushRequest{Key: key}

	switch v := value.(type) {
	case string:
		req.Value = &glist.GListRPushRequest_StringValue{StringValue: v}
	case int32:
		req.Value = &glist.GListRPushRequest_Int32Value{Int32Value: v}
	case int64:
		req.Value = &glist.GListRPushRequest_Int64Value{Int64Value: v}
	case float32:
		req.Value = &glist.GListRPushRequest_Float32Value{Float32Value: v}
	case float64:
		req.Value = &glist.GListRPushRequest_Float64Value{Float64Value: v}
	case bool:
		req.Value = &glist.GListRPushRequest_BoolValue{BoolValue: v}
	case []byte:
		req.Value = &glist.GListRPushRequest_BytesValue{BytesValue: v}
	default:
		return errors.New("unknown value type")
	}

	rpush, err := client.RPush(context.Background(), req)
	if err != nil {
		return errors.New("client RPush failed: " + err.Error())
	}
	if !rpush.Ok {
		return errors.New("RPush failed")
	}
	return nil
}

func (c *Client) RPushs(key string, values []interface{}) error {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListRPushsRequest{Key: key}
	// Convert the list of values to gRPC-compatible format
	for _, value := range values {
		//print(value)
		switch v := value.(type) {
		case string:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_StringValue{StringValue: v}})
		case int32:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_Int32Value{Int32Value: v}})
		case int64:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_Int64Value{Int64Value: v}})
		case float32:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_Float32Value{Float32Value: v}})
		case float64:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_Float64Value{Float64Value: v}})
		case bool:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_BoolValue{BoolValue: v}})
		case []byte:
			req.Values = append(req.Values, &glist.Value{Value: &glist.Value_BytesValue{BytesValue: v}})
		default:
			return errors.New("unknown value type")
		}
	}
	print(111)
	_, err = client.RPushs(context.Background(), req)
	if err != nil {
		return errors.New("client RPushs failed: " + err.Error())
	}
	print("222")

	return nil
}

func (c *Client) LPop(key string) (interface{}, error) {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListLPopRequest{Key: key}
	resp, err := client.LPop(context.Background(), req)
	if err != nil {
		return nil, errors.New("client LPop failed: " + err.Error())
	}

	// Convert the response to the appropriate type based on the gRPC response
	switch v := resp.Value.(type) {
	case *glist.GListLPopResponse_StringValue:
		return v.StringValue, nil
	case *glist.GListLPopResponse_Int32Value:
		return v.Int32Value, nil
	case *glist.GListLPopResponse_Int64Value:
		return v.Int64Value, nil
	case *glist.GListLPopResponse_Float32Value:
		return v.Float32Value, nil
	case *glist.GListLPopResponse_Float64Value:
		return v.Float64Value, nil
	case *glist.GListLPopResponse_BoolValue:
		return v.BoolValue, nil
	case *glist.GListLPopResponse_BytesValue:
		return v.BytesValue, nil
	default:
		return nil, errors.New("unknown value type")
	}
}

func (c *Client) RPop(key string) (interface{}, error) {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListRPopRequest{Key: key}
	resp, err := client.RPop(context.Background(), req)
	if err != nil {
		return nil, errors.New("client RPop failed: " + err.Error())
	}

	// Convert the response to the appropriate type based on the gRPC response
	switch v := resp.Value.(type) {
	case *glist.GListRPopResponse_StringValue:
		return v.StringValue, nil
	case *glist.GListRPopResponse_Int32Value:
		return v.Int32Value, nil
	case *glist.GListRPopResponse_Int64Value:
		return v.Int64Value, nil
	case *glist.GListRPopResponse_Float32Value:
		return v.Float32Value, nil
	case *glist.GListRPopResponse_Float64Value:
		return v.Float64Value, nil
	case *glist.GListRPopResponse_BoolValue:
		return v.BoolValue, nil
	case *glist.GListRPopResponse_BytesValue:
		return v.BytesValue, nil
	default:
		return nil, errors.New("unknown value type")
	}
}

func (c *Client) LRange(key string, start, stop int) ([]interface{}, error) {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListLRangeRequest{
		Key:   key,
		Start: int32(start),
		Stop:  int32(stop),
	}

	resp, err := client.LRange(context.Background(), req)
	if err != nil {
		return nil, errors.New("client LRange failed: " + err.Error())
	}

	// Convert the response to a slice of interface{} based on the gRPC response
	var values []interface{}
	for _, val := range resp.Values {
		switch v := val.Value.(type) {
		case *glist.Value_StringValue:
			values = append(values, v.StringValue)
		case *glist.Value_Int32Value:
			values = append(values, v.Int32Value)
		case *glist.Value_Int64Value:
			values = append(values, v.Int64Value)
		case *glist.Value_Float32Value:
			values = append(values, v.Float32Value)
		case *glist.Value_Float64Value:
			values = append(values, v.Float64Value)
		case *glist.Value_BoolValue:
			values = append(values, v.BoolValue)
		case *glist.Value_BytesValue:
			values = append(values, v.BytesValue)
		default:
			return nil, errors.New("unknown value type")
		}
	}

	return values, nil
}

func (c *Client) LLen(key string) (int32, error) {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return 0, errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListLLenRequest{Key: key}
	resp, err := client.LLen(context.Background(), req)
	if err != nil {
		return 0, errors.New("client LLen failed: " + err.Error())
	}

	return resp.Length, nil
}

func (c *Client) LRem(key string, count int32, value interface{}) error {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListLRemRequest{
		Key:   key,
		Count: count,
	}

	switch v := value.(type) {
	case string:
		req.Value = &glist.GListLRemRequest_StringValue{StringValue: v}
	case int32:
		req.Value = &glist.GListLRemRequest_Int32Value{Int32Value: v}
	case int64:
		req.Value = &glist.GListLRemRequest_Int64Value{Int64Value: v}
	case float32:
		req.Value = &glist.GListLRemRequest_Float32Value{Float32Value: v}
	case float64:
		req.Value = &glist.GListLRemRequest_Float64Value{Float64Value: v}
	case bool:
		req.Value = &glist.GListLRemRequest_BoolValue{BoolValue: v}
	case []byte:
		req.Value = &glist.GListLRemRequest_BytesValue{BytesValue: v}
	default:
		return errors.New("unknown value type")
	}

	_, err = client.LRem(context.Background(), req)
	if err != nil {
		return errors.New("client LRem failed: " + err.Error())
	}
	return nil
}

func (c *Client) LIndex(key string, index int) (interface{}, error) {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListLIndexRequest{
		Key:   key,
		Index: int32(index),
	}

	resp, err := client.LIndex(context.Background(), req)
	if err != nil {
		return nil, errors.New("client LIndex failed: " + err.Error())
	}

	switch v := resp.Value.(type) {
	case *glist.GListLIndexResponse_StringValue:
		return v.StringValue, nil
	case *glist.GListLIndexResponse_Int32Value:
		return v.Int32Value, nil
	case *glist.GListLIndexResponse_Int64Value:
		return v.Int64Value, nil
	case *glist.GListLIndexResponse_Float32Value:
		return v.Float32Value, nil
	case *glist.GListLIndexResponse_Float64Value:
		return v.Float64Value, nil
	case *glist.GListLIndexResponse_BoolValue:
		return v.BoolValue, nil
	case *glist.GListLIndexResponse_BytesValue:
		return v.BytesValue, nil
	default:
		return nil, errors.New("unknown value type")
	}
}

func (c *Client) LSet(key string, index int, value interface{}) error {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListLSetRequest{
		Key:   key,
		Index: int32(index),
	}

	switch v := value.(type) {
	case string:
		req.Value = &glist.GListLSetRequest_StringValue{StringValue: v}
	case int32:
		req.Value = &glist.GListLSetRequest_Int32Value{Int32Value: v}
	case int64:
		req.Value = &glist.GListLSetRequest_Int64Value{Int64Value: v}
	case float32:
		req.Value = &glist.GListLSetRequest_Float32Value{Float32Value: v}
	case float64:
		req.Value = &glist.GListLSetRequest_Float64Value{Float64Value: v}
	case bool:
		req.Value = &glist.GListLSetRequest_BoolValue{BoolValue: v}
	case []byte:
		req.Value = &glist.GListLSetRequest_BytesValue{BytesValue: v}
	default:
		return errors.New("unknown value type")
	}

	resp, err := client.LSet(context.Background(), req)
	if err != nil {
		return errors.New("client LSet failed: " + err.Error())
	}

	if !resp.Ok {
		return errors.New("LSet failed")
	}
	return nil
}

func (c *Client) LTrim(key string, start, stop int) error {
	client, err := newListGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}

	req := &glist.GListLTrimRequest{
		Key:   key,
		Start: int32(start),
		Stop:  int32(stop),
	}

	resp, err := client.LTrim(context.Background(), req)
	if err != nil {
		return errors.New("client LTrim failed: " + err.Error())
	}

	if !resp.Ok {
		return errors.New("LTrim failed")
	}
	return nil
}
