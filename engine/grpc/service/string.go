package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/proto/gstring"
	"github.com/ByteStorage/FlyDB/structure"
	"time"
)

type StringService interface {
	Base
	gstring.GStringServiceServer
}

type str struct {
	dbs *structure.StringStructure
	gstring.GStringServiceServer
}

func (s *str) CloseDb() error {
	return s.dbs.Stop()
}

func NewStringService(options config.Options) (StringService, error) {
	stringStructure, err := structure.NewStringStructure(options)
	if err != nil {
		return nil, err
	}
	return &str{
		dbs: stringStructure,
	}, nil
}

func (s *str) NewFlyDBService(ctx context.Context, req *gstring.FlyDBOption) (*gstring.NewFlyDBResponse, error) {
	option := config.DefaultOptions
	if req.DirPath != "" {
		option.DirPath = req.DirPath
	}
	if req.DataFileSize != 0 {
		option.DataFileSize = req.DataFileSize
	}
	if req.SyncWrite {
		option.SyncWrite = req.SyncWrite
	}
	fmt.Println("new flydb option: ", req.DirPath)
	dbs, err := structure.NewStringStructure(option)
	if err != nil {
		return &gstring.NewFlyDBResponse{
			ResponseMsg: err.Error(),
		}, nil
	}
	s.dbs = dbs
	return &gstring.NewFlyDBResponse{
		ResponseMsg: "start success!",
	}, nil
}

// Put is a grpc s for put
func (s *str) Put(ctx context.Context, req *gstring.SetRequest) (*gstring.SetResponse, error) {
	fmt.Println("receive put request: key: ", req.Key, " value: ", req.GetValue(), " duration: ", time.Duration(req.Expire))
	var err error
	switch req.Value.(type) {
	case *gstring.SetRequest_StringValue:
		err = s.dbs.Set(req.Key, req.GetStringValue(), time.Duration(req.Expire))
	case *gstring.SetRequest_Int32Value:
		err = s.dbs.Set(req.Key, req.GetInt32Value(), time.Duration(req.Expire))
	case *gstring.SetRequest_Int64Value:
		err = s.dbs.Set(req.Key, req.GetInt64Value(), time.Duration(req.Expire))
	case *gstring.SetRequest_Float32Value:
		err = s.dbs.Set(req.Key, req.GetFloat32Value(), time.Duration(req.Expire))
	case *gstring.SetRequest_Float64Value:
		err = s.dbs.Set(req.Key, req.GetFloat64Value(), time.Duration(req.Expire))
	case *gstring.SetRequest_BoolValue:
		err = s.dbs.Set(req.Key, req.GetBoolValue(), time.Duration(req.Expire))
	case *gstring.SetRequest_BytesValue:
		err = s.dbs.Set(req.Key, req.GetBytesValue(), time.Duration(req.Expire))
	default:
		err = fmt.Errorf("unknown value type")
	}
	if err != nil {
		return &gstring.SetResponse{}, err
	}
	return &gstring.SetResponse{Ok: true}, nil
}

// Get is a grpc s for get
func (s *str) Get(ctx context.Context, req *gstring.GetRequest) (*gstring.GetResponse, error) {
	value, err := s.dbs.Get(req.Key)
	if err != nil {
		return &gstring.GetResponse{}, err
	}
	resp := &gstring.GetResponse{}
	switch v := value.(type) {
	case string:
		resp.Value = &gstring.GetResponse_StringValue{StringValue: v}
	case int32:
		resp.Value = &gstring.GetResponse_Int32Value{Int32Value: v}
	case int64:
		resp.Value = &gstring.GetResponse_Int64Value{Int64Value: v}
	case float32:
		resp.Value = &gstring.GetResponse_Float32Value{Float32Value: v}
	case float64:
		resp.Value = &gstring.GetResponse_Float64Value{Float64Value: v}
	case bool:
		resp.Value = &gstring.GetResponse_BoolValue{BoolValue: v}
	case []byte:
		resp.Value = &gstring.GetResponse_BytesValue{BytesValue: v}
	}
	return resp, nil
}

// Del is a grpc s for del
func (s *str) Del(ctx context.Context, req *gstring.DelRequest) (*gstring.DelResponse, error) {
	err := s.dbs.Del(req.Key)
	if err != nil {
		return &gstring.DelResponse{}, err
	}
	return &gstring.DelResponse{Ok: true}, nil
}

func (s *str) Type(ctx context.Context, req *gstring.TypeRequest) (*gstring.TypeResponse, error) {
	strtype, err := s.dbs.Type(req.Key)
	if err != nil {
		return &gstring.TypeResponse{}, err
	}
	return &gstring.TypeResponse{Type: strtype}, nil
}

func (s *str) StrLen(ctx context.Context, req *gstring.StrLenRequest) (*gstring.StrLenResponse, error) {
	length, err := s.dbs.StrLen(req.Key)
	if err != nil {
		return &gstring.StrLenResponse{}, err
	}
	return &gstring.StrLenResponse{
		Length: int32(length),
	}, nil
}

func (s *str) GetSet(ctx context.Context, req *gstring.GetSetRequest) (*gstring.GetSetResponse, error) {
	var previousValue interface{}
	var err error

	switch v := req.Value.(type) {
	case *gstring.GetSetRequest_StringValue:
		previousValue, err = s.dbs.GetSet(req.Key, v.StringValue, time.Duration(req.Expire))
	case *gstring.GetSetRequest_Int32Value:
		previousValue, err = s.dbs.GetSet(req.Key, v.Int32Value, time.Duration(req.Expire))
	case *gstring.GetSetRequest_Int64Value:
		previousValue, err = s.dbs.GetSet(req.Key, v.Int64Value, time.Duration(req.Expire))
	case *gstring.GetSetRequest_Float32Value:
		previousValue, err = s.dbs.GetSet(req.Key, v.Float32Value, time.Duration(req.Expire))
	case *gstring.GetSetRequest_Float64Value:
		previousValue, err = s.dbs.GetSet(req.Key, v.Float64Value, time.Duration(req.Expire))
	case *gstring.GetSetRequest_BoolValue:
		previousValue, err = s.dbs.GetSet(req.Key, v.BoolValue, time.Duration(req.Expire))
	case *gstring.GetSetRequest_BytesValue:
		previousValue, err = s.dbs.GetSet(req.Key, v.BytesValue, time.Duration(req.Expire))
	default:
		return nil, fmt.Errorf("unknown value type")
	}
	if err != nil {
		return &gstring.GetSetResponse{}, err
	}
	resp := &gstring.GetSetResponse{}
	switch v := previousValue.(type) {
	case string:
		resp.Value = &gstring.GetSetResponse_StringValue{StringValue: v}
	case int32:
		resp.Value = &gstring.GetSetResponse_Int32Value{Int32Value: v}
	case int64:
		resp.Value = &gstring.GetSetResponse_Int64Value{Int64Value: v}
	case float32:
		resp.Value = &gstring.GetSetResponse_Float32Value{Float32Value: v}
	case float64:
		resp.Value = &gstring.GetSetResponse_Float64Value{Float64Value: v}
	case bool:
		resp.Value = &gstring.GetSetResponse_BoolValue{BoolValue: v}
	case []byte:
		resp.Value = &gstring.GetSetResponse_BytesValue{BytesValue: v}
	}

	return resp, nil
}

func (s *str) Append(ctx context.Context, req *gstring.AppendRequest) (*gstring.AppendResponse, error) {
	err := s.dbs.Append(req.Key, req.Value, time.Duration(req.Expire))
	if err != nil {
		return &gstring.AppendResponse{}, err
	}
	return &gstring.AppendResponse{Ok: true}, nil
}

func (s *str) Incr(ctx context.Context, req *gstring.IncrRequest) (*gstring.IncrResponse, error) {
	err := s.dbs.Incr(req.Key, time.Duration(req.Expire))
	if err != nil {
		return &gstring.IncrResponse{}, err
	}
	return &gstring.IncrResponse{Ok: true}, nil
}

func (s *str) IncrBy(ctx context.Context, req *gstring.IncrByRequest) (*gstring.IncrByResponse, error) {
	err := s.dbs.IncrBy(req.Key, int(req.Amount), time.Duration(req.Expire))
	if err != nil {
		return &gstring.IncrByResponse{}, err
	}
	return &gstring.IncrByResponse{Ok: true}, nil
}

func (s *str) IncrByFloat(ctx context.Context, req *gstring.IncrByFloatRequest) (*gstring.IncrByFloatResponse, error) {
	err := s.dbs.IncrByFloat(req.Key, req.Amount, time.Duration(req.Expire))
	if err != nil {
		return &gstring.IncrByFloatResponse{Ok: true}, err
	}
	return &gstring.IncrByFloatResponse{Ok: true}, nil
}

func (s *str) Decr(ctx context.Context, req *gstring.DecrRequest) (*gstring.DecrResponse, error) {
	err := s.dbs.Decr(req.Key, time.Duration(req.Expire))
	if err != nil {
		return &gstring.DecrResponse{}, err
	}
	return &gstring.DecrResponse{Ok: true}, nil
}

func (s *str) DecrBy(ctx context.Context, req *gstring.DecrByRequest) (*gstring.DecrByResponse, error) {
	err := s.dbs.DecrBy(req.Key, int(req.Amount), time.Duration(req.Expire))
	if err != nil {
		return &gstring.DecrByResponse{}, err
	}
	return &gstring.DecrByResponse{Ok: true}, nil
}

func (s *str) Exists(ctx context.Context, req *gstring.ExistsRequest) (*gstring.ExistsResponse, error) {
	exists, err := s.dbs.Exists(req.Key)
	if err != nil {
		return nil, err
	}

	return &gstring.ExistsResponse{Exists: exists}, nil
}

func (s *str) Expire(ctx context.Context, req *gstring.ExpireRequest) (*gstring.ExpireResponse, error) {
	err := s.dbs.Expire(req.Key, time.Duration(req.Expire)*time.Second)
	if err != nil {
		return &gstring.ExpireResponse{}, err
	}
	return &gstring.ExpireResponse{Ok: true}, nil
}

func (s *str) Persist(ctx context.Context, req *gstring.PersistRequest) (*gstring.PersistResponse, error) {
	err := s.dbs.Persist(req.Key)
	if err != nil {
		return nil, err
	}

	return &gstring.PersistResponse{Ok: true}, nil
}

func (s *str) MGet(ctx context.Context, req *gstring.MGetRequest) (*gstring.MGetResponse, error) {
	values, err := s.dbs.MGet(req.Keys...)
	if err != nil {
		return &gstring.MGetResponse{}, err
	}
	resp := &gstring.MGetResponse{}

	for _, value := range values {
		switch v := value.(type) {
		case string:
			resp.Values = append(resp.Values, &gstring.MGetValue{Value: &gstring.MGetValue_StringValue{StringValue: v}})
		case int32:
			resp.Values = append(resp.Values, &gstring.MGetValue{Value: &gstring.MGetValue_Int32Value{Int32Value: v}})
		case int64:
			resp.Values = append(resp.Values, &gstring.MGetValue{Value: &gstring.MGetValue_Int64Value{Int64Value: v}})
		case float32:
			resp.Values = append(resp.Values, &gstring.MGetValue{Value: &gstring.MGetValue_Float32Value{Float32Value: v}})
		case float64:
			resp.Values = append(resp.Values, &gstring.MGetValue{Value: &gstring.MGetValue_Float64Value{Float64Value: v}})
		case bool:
			resp.Values = append(resp.Values, &gstring.MGetValue{Value: &gstring.MGetValue_BoolValue{BoolValue: v}})
		case []byte:
			resp.Values = append(resp.Values, &gstring.MGetValue{Value: &gstring.MGetValue_BytesValue{BytesValue: v}})
		}
	}

	return resp, nil
}

func (s *str) MSet(ctx context.Context, req *gstring.MSetRequest) (*gstring.MSetResponse, error) {
	// Extract the key-value pairs from the request and populate the `values` slice.
	var values []interface{}
	for _, keyValue := range req.GetPairs() {
		switch v := keyValue.Value.(type) {
		case *gstring.MSetRequest_KeyValue_StringValue:
			values = append(values, keyValue.GetKey(), v.StringValue)
		case *gstring.MSetRequest_KeyValue_Int32Value:
			values = append(values, keyValue.GetKey(), v.Int32Value)
		case *gstring.MSetRequest_KeyValue_Int64Value:
			values = append(values, keyValue.GetKey(), v.Int64Value)
		case *gstring.MSetRequest_KeyValue_Float32Value:
			values = append(values, keyValue.GetKey(), v.Float32Value)
		case *gstring.MSetRequest_KeyValue_Float64Value:
			values = append(values, keyValue.GetKey(), v.Float64Value)
		case *gstring.MSetRequest_KeyValue_BoolValue:
			values = append(values, keyValue.GetKey(), v.BoolValue)
		case *gstring.MSetRequest_KeyValue_BytesValue:
			values = append(values, keyValue.GetKey(), v.BytesValue)
		default:
			return nil, errors.New("unsupported value type")
		}
	}
	//print(values)
	// Use the `MSet` method of the store to set the key-value pairs.
	err := s.dbs.MSet(values...)
	if err != nil {
		return nil, err
	}

	// Create the response indicating success.
	response := &gstring.MSetResponse{
		Success: true,
	}

	return response, nil
}

func (s *str) MSetNX(ctx context.Context, req *gstring.MSetNXRequest) (*gstring.MSetNXResponse, error) {
	// Extract the key-value pairs from the request and populate the `values` slice.
	var values []interface{}
	for _, keyValue := range req.GetPairs() {
		switch v := keyValue.Value.(type) {
		case *gstring.MSetNXRequest_KeyValue_StringValue:
			values = append(values, keyValue.GetKey(), v.StringValue)
		case *gstring.MSetNXRequest_KeyValue_Int32Value:
			values = append(values, keyValue.GetKey(), v.Int32Value)
		case *gstring.MSetNXRequest_KeyValue_Int64Value:
			values = append(values, keyValue.GetKey(), v.Int64Value)
		case *gstring.MSetNXRequest_KeyValue_Float32Value:
			values = append(values, keyValue.GetKey(), v.Float32Value)
		case *gstring.MSetNXRequest_KeyValue_Float64Value:
			values = append(values, keyValue.GetKey(), v.Float64Value)
		case *gstring.MSetNXRequest_KeyValue_BoolValue:
			values = append(values, keyValue.GetKey(), v.BoolValue)
		case *gstring.MSetNXRequest_KeyValue_BytesValue:
			values = append(values, keyValue.GetKey(), v.BytesValue)
		default:
			return nil, errors.New("unsupported value type")
		}
	}
	//print(values)
	// Use the `MSet` method of the store to set the key-value pairs.
	exists, err := s.dbs.MSetNX(values...)
	if err != nil {
		return nil, err
	}

	// Create the response indicating success.
	response := &gstring.MSetNXResponse{
		Success: exists,
	}

	return response, nil
}
