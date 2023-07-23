package service

import (
	"context"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/proto/glist"
	"github.com/ByteStorage/FlyDB/structure"
)

type ListService interface {
	Base
	glist.GListServiceServer
}

type list struct {
	dbs *structure.ListStructure
	glist.GListServiceServer
}

func (l *list) CloseDb() error {
	return l.dbs.Stop()
}

func NewListService(options config.Options) (ListService, error) {
	listStructure, err := structure.NewListStructure(options)
	if err != nil {
		return nil, err
	}
	return &list{dbs: listStructure}, nil
}

func (l *list) LPush(ctx context.Context, req *glist.GListLPushRequest) (*glist.GListLPushResponse, error) {
	// Implement the logic for the LPush gRPC method
	// ...	\
	var err error
	switch req.Value.(type) {
	case *glist.GListLPushRequest_StringValue:
		err = l.dbs.LPush(req.Key, req.GetStringValue())
	case *glist.GListLPushRequest_Int64Value:
		err = l.dbs.LPush(req.Key, req.GetInt64Value())
	case *glist.GListLPushRequest_Float32Value:
		err = l.dbs.LPush(req.Key, req.GetFloat32Value())
	case *glist.GListLPushRequest_Float64Value:
		err = l.dbs.LPush(req.Key, req.GetFloat64Value())
	case *glist.GListLPushRequest_BoolValue:
		err = l.dbs.LPush(req.Key, req.GetBoolValue())
	case *glist.GListLPushRequest_BytesValue:
		err = l.dbs.LPush(req.Key, req.GetBytesValue())
	default:
		err = fmt.Errorf("unknown value type")
	}
	if err != nil {
		return &glist.GListLPushResponse{}, err
	}
	return &glist.GListLPushResponse{Ok: true}, nil
}

func (l *list) LPushs(ctx context.Context, req *glist.GListLPushsRequest) (*glist.GListLPushsResponse, error) {
	values := make([]interface{}, len(req.Values))

	for i, value := range req.Values {
		switch v := value.Value.(type) {
		case *glist.Value_StringValue:
			values[i] = v.StringValue
		case *glist.Value_Int32Value:
			values[i] = v.Int32Value
		case *glist.Value_Int64Value:
			values[i] = v.Int64Value
		case *glist.Value_Float32Value:
			values[i] = v.Float32Value
		case *glist.Value_Float64Value:
			values[i] = v.Float64Value
		case *glist.Value_BoolValue:
			values[i] = v.BoolValue
		case *glist.Value_BytesValue:
			values[i] = v.BytesValue
		default:
			return nil, fmt.Errorf("unsupported value type")
		}
	}

	err := l.dbs.LPushs(req.Key, values...)
	if err != nil {
		return &glist.GListLPushsResponse{}, err
	}

	return &glist.GListLPushsResponse{Ok: true}, nil
}

func (l *list) RPush(ctx context.Context, req *glist.GListRPushRequest) (*glist.GListRPushResponse, error) {
	// Implement the logic for the RPush gRPC method
	// ...
	var err error
	switch req.Value.(type) {
	case *glist.GListRPushRequest_StringValue:
		err = l.dbs.RPush(req.Key, req.GetStringValue())
	case *glist.GListRPushRequest_Int64Value:
		err = l.dbs.RPush(req.Key, req.GetInt64Value())
	case *glist.GListRPushRequest_Float32Value:
		err = l.dbs.RPush(req.Key, req.GetFloat32Value())
	case *glist.GListRPushRequest_Float64Value:
		err = l.dbs.RPush(req.Key, req.GetFloat64Value())
	case *glist.GListRPushRequest_BoolValue:
		err = l.dbs.RPush(req.Key, req.GetBoolValue())
	case *glist.GListRPushRequest_BytesValue:
		err = l.dbs.RPush(req.Key, req.GetBytesValue())
	default:
		err = fmt.Errorf("unknown value type")
	}
	if err != nil {
		return &glist.GListRPushResponse{}, err
	}
	return &glist.GListRPushResponse{Ok: true}, nil
}

func (l *list) RPushs(ctx context.Context, req *glist.GListRPushsRequest) (*glist.GListRPushsResponse, error) {
	values := make([]interface{}, len(req.Values))

	for i, value := range req.Values {
		switch v := value.Value.(type) {
		case *glist.Value_StringValue:
			values[i] = v.StringValue
		case *glist.Value_Int32Value:
			values[i] = v.Int32Value
		case *glist.Value_Int64Value:
			values[i] = v.Int64Value
		case *glist.Value_Float32Value:
			values[i] = v.Float32Value
		case *glist.Value_Float64Value:
			values[i] = v.Float64Value
		case *glist.Value_BoolValue:
			values[i] = v.BoolValue
		case *glist.Value_BytesValue:
			values[i] = v.BytesValue
		default:
			return nil, fmt.Errorf("unsupported value type")
		}
	}

	err := l.dbs.RPushs(req.Key, values...)
	if err != nil {
		return &glist.GListRPushsResponse{}, err
	}

	return &glist.GListRPushsResponse{Ok: true}, nil
}

func (l *list) LPop(ctx context.Context, req *glist.GListLPopRequest) (*glist.GListLPopResponse, error) {
	// Implement the logic for the LPop gRPC method
	// ...
	value, err := l.dbs.LPop(req.Key)
	if err != nil {
		return &glist.GListLPopResponse{}, err
	}
	resp := &glist.GListLPopResponse{}
	switch v := value.(type) {
	case string:
		resp.Value = &glist.GListLPopResponse_StringValue{StringValue: v}
	case int32:
		resp.Value = &glist.GListLPopResponse_Int32Value{Int32Value: v}
	case int64:
		resp.Value = &glist.GListLPopResponse_Int64Value{Int64Value: v}
	case float32:
		resp.Value = &glist.GListLPopResponse_Float32Value{Float32Value: v}
	case float64:
		resp.Value = &glist.GListLPopResponse_Float64Value{Float64Value: v}
	case bool:
		resp.Value = &glist.GListLPopResponse_BoolValue{BoolValue: v}
	case []byte:
		resp.Value = &glist.GListLPopResponse_BytesValue{BytesValue: v}
	}
	return resp, nil
}

func (l *list) RPop(ctx context.Context, req *glist.GListRPopRequest) (*glist.GListRPopResponse, error) {
	// Implement the logic for the RPop gRPC method
	// ...
	value, err := l.dbs.RPop(req.Key)
	if err != nil {
		return &glist.GListRPopResponse{}, err
	}
	resp := &glist.GListRPopResponse{}
	switch v := value.(type) {
	case string:
		resp.Value = &glist.GListRPopResponse_StringValue{StringValue: v}
	case int32:
		resp.Value = &glist.GListRPopResponse_Int32Value{Int32Value: v}
	case int64:
		resp.Value = &glist.GListRPopResponse_Int64Value{Int64Value: v}
	case float32:
		resp.Value = &glist.GListRPopResponse_Float32Value{Float32Value: v}
	case float64:
		resp.Value = &glist.GListRPopResponse_Float64Value{Float64Value: v}
	case bool:
		resp.Value = &glist.GListRPopResponse_BoolValue{BoolValue: v}
	case []byte:
		resp.Value = &glist.GListRPopResponse_BytesValue{BytesValue: v}
	}
	return resp, nil
}

func (l *list) LRange(ctx context.Context, req *glist.GListLRangeRequest) (*glist.GListLRangeResponse, error) {
	// Implement the logic for the LRange gRPC method
	listValues, err := l.dbs.LRange(req.Key, int(req.Start), int(req.Stop))
	if err != nil {
		return nil, err
	}
	// Convert the list of values to the protobuf response format
	resp := &glist.GListLRangeResponse{}
	for _, val := range listValues {
		switch v := val.(type) {
		case string:
			resp.Values = append(resp.Values, &glist.Value{Value: &glist.Value_StringValue{StringValue: v}})
		case int32:
			resp.Values = append(resp.Values, &glist.Value{Value: &glist.Value_Int32Value{Int32Value: v}})
		case int64:
			resp.Values = append(resp.Values, &glist.Value{Value: &glist.Value_Int64Value{Int64Value: v}})
		case float32:
			resp.Values = append(resp.Values, &glist.Value{Value: &glist.Value_Float32Value{Float32Value: v}})
		case float64:
			resp.Values = append(resp.Values, &glist.Value{Value: &glist.Value_Float64Value{Float64Value: v}})
		case bool:
			resp.Values = append(resp.Values, &glist.Value{Value: &glist.Value_BoolValue{BoolValue: v}})
		case []byte:
			resp.Values = append(resp.Values, &glist.Value{Value: &glist.Value_BytesValue{BytesValue: v}})
		}

	}
	return resp, nil
}

func (l *list) LLen(ctx context.Context, req *glist.GListLLenRequest) (*glist.GListLLenResponse, error) {
	// Implement the logic for the LLen gRPC method
	length, err := l.dbs.LLen(req.Key)
	if err != nil {
		return nil, err
	}

	return &glist.GListLLenResponse{Length: int32(length)}, nil
}

func (l *list) LRem(ctx context.Context, req *glist.GListLRemRequest) (*glist.GListLRemResponse, error) {
	err := l.dbs.LRem(req.Key, int(req.Count), req.Value)
	if err != nil {
		return &glist.GListLRemResponse{}, err
	}
	return &glist.GListLRemResponse{Ok: true}, nil
}

func (l *list) LIndex(ctx context.Context, req *glist.GListLIndexRequest) (*glist.GListLIndexResponse, error) {
	value, err := l.dbs.LIndex(req.Key, int(req.Index))
	if err != nil {
		return nil, err
	}
	resp := &glist.GListLIndexResponse{}
	switch v := value.(type) {
	case string:
		resp.Value = &glist.GListLIndexResponse_StringValue{StringValue: v}
	case int32:
		resp.Value = &glist.GListLIndexResponse_Int32Value{Int32Value: v}
	case int64:
		resp.Value = &glist.GListLIndexResponse_Int64Value{Int64Value: v}
	case float32:
		resp.Value = &glist.GListLIndexResponse_Float32Value{Float32Value: v}
	case float64:
		resp.Value = &glist.GListLIndexResponse_Float64Value{Float64Value: v}
	case bool:
		resp.Value = &glist.GListLIndexResponse_BoolValue{BoolValue: v}
	case []byte:
		resp.Value = &glist.GListLIndexResponse_BytesValue{BytesValue: v}
	}
	return resp, nil
}

func (l *list) LSet(ctx context.Context, req *glist.GListLSetRequest) (*glist.GListLSetResponse, error) {
	var err error
	switch req.Value.(type) {
	case *glist.GListLSetRequest_StringValue:
		err = l.dbs.LSet(req.Key, int(req.Index), req.GetStringValue())
	case *glist.GListLSetRequest_Int64Value:
		err = l.dbs.LSet(req.Key, int(req.Index), req.GetInt64Value())
	case *glist.GListLSetRequest_Float32Value:
		err = l.dbs.LSet(req.Key, int(req.Index), req.GetFloat32Value())
	case *glist.GListLSetRequest_Float64Value:
		err = l.dbs.LSet(req.Key, int(req.Index), req.GetFloat64Value())
	case *glist.GListLSetRequest_BoolValue:
		err = l.dbs.LSet(req.Key, int(req.Index), req.GetBoolValue())
	case *glist.GListLSetRequest_BytesValue:
		err = l.dbs.LSet(req.Key, int(req.Index), req.GetBytesValue())
	default:
		err = fmt.Errorf("unknown value type")
	}
	if err != nil {
		return &glist.GListLSetResponse{}, err
	}
	return &glist.GListLSetResponse{Ok: true}, nil
}

func (l *list) LTrim(ctx context.Context, req *glist.GListLTrimRequest) (*glist.GListLTrimResponse, error) {
	err := l.dbs.LTrim(req.Key, int(req.Start), int(req.Stop))
	if err != nil {
		return &glist.GListLTrimResponse{}, err
	}

	return &glist.GListLTrimResponse{Ok: true}, nil
}
