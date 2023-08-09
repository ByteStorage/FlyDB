package service

import (
	"context"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/proto/ghash"
	"github.com/ByteStorage/FlyDB/structure"
)

type HashService interface {
	Base
	ghash.GHashServiceServer
}

type hash struct {
	dbh *structure.HashStructure
	ghash.GHashServiceServer
}

func NewHashService(options config.Options) (HashService, error) {
	dbh, err := structure.NewHashStructure(options)
	if err != nil {
		return nil, err
	}
	return &hash{dbh: dbh}, nil
}

func (s *hash) CloseDb() error {
	return s.dbh.Stop()
}

// HSet is a grpc s for put
func (s *hash) HSet(ctx context.Context, req *ghash.GHashSetRequest) (*ghash.GHashSetResponse, error) {
	fmt.Println("receive put request: key: ", req.Key, " field: ", req.GetField(), " value: ", req.GetValue())
	var err error
	result, err := setValue(s, req.Key, req.Field, req.Ttl, req)
	if err != nil {
		return &ghash.GHashSetResponse{Ok: result}, err
	}
	return &ghash.GHashSetResponse{Ok: result}, nil

}

// HGet is a grpc s for get
func (s *hash) HGet(ctx context.Context, req *ghash.GHashGetRequest) (*ghash.GHashGetResponse, error) {
	value, err := s.dbh.HGet(req.Key, req.Field)
	if err != nil {
		return &ghash.GHashGetResponse{}, err
	}
	resp := &ghash.GHashGetResponse{}
	switch v := value.(type) {
	case string:
		resp.Value = &ghash.GHashGetResponse_StringValue{StringValue: v}
	case int32:
		resp.Value = &ghash.GHashGetResponse_Int32Value{Int32Value: v}
	case int64:
		resp.Value = &ghash.GHashGetResponse_Int64Value{Int64Value: v}
	case float32:
		resp.Value = &ghash.GHashGetResponse_Float32Value{Float32Value: v}
	case float64:
		resp.Value = &ghash.GHashGetResponse_Float64Value{Float64Value: v}
	case bool:
		resp.Value = &ghash.GHashGetResponse_BoolValue{BoolValue: v}
	case []byte:
		resp.Value = &ghash.GHashGetResponse_BytesValue{BytesValue: v}
	}
	return resp, nil
}

// HDel is a grpc s for del
func (s *hash) HDel(ctx context.Context, req *ghash.GHashDelRequest) (*ghash.GHashDelResponse, error) {
	ok, err := s.dbh.HDel(req.Key, req.Field)
	if err != nil {
		return &ghash.GHashDelResponse{Ok: ok}, err
	}
	return &ghash.GHashDelResponse{Ok: ok}, err
}

func setValue(s *hash, key string, field interface{}, ttl int64, r *ghash.GHashSetRequest) (bool, error) {
	switch r.Value.(type) {
	case *ghash.GHashSetRequest_StringValue:
		ok, err := s.dbh.HSet(key, field, r.GetStringValue(), ttl)
		return ok, err
	case *ghash.GHashSetRequest_Int32Value:
		ok, err := s.dbh.HSet(key, field, r.GetInt32Value(), ttl)
		return ok, err
	case *ghash.GHashSetRequest_Int64Value:
		ok, err := s.dbh.HSet(key, field, r.GetInt64Value(), ttl)
		return ok, err
	case *ghash.GHashSetRequest_Float32Value:
		ok, err := s.dbh.HSet(key, field, r.GetFloat32Value(), ttl)
		return ok, err
	case *ghash.GHashSetRequest_Float64Value:
		ok, err := s.dbh.HSet(key, field, r.GetFloat64Value(), ttl)
		return ok, err
	case *ghash.GHashSetRequest_BoolValue:
		ok, err := s.dbh.HSet(key, field, r.GetBoolValue(), ttl)
		return ok, err
	case *ghash.GHashSetRequest_BytesValue:
		ok, err := s.dbh.HSet(key, field, r.GetBytesValue(), ttl)
		return ok, err
	default:
		return false, fmt.Errorf("unknown value type")
	}
}

func (s *hash) HExists(ctx context.Context, req *ghash.GHashExistsRequest) (*ghash.GHashExistsResponse, error) {
	exists, err := s.dbh.HExists(req.Key, req.Field)
	if err != nil {
		return &ghash.GHashExistsResponse{Ok: exists}, err
	}
	return &ghash.GHashExistsResponse{Ok: exists}, nil
}

func (s *hash) HLen(ctx context.Context, req *ghash.GHashLenRequest) (*ghash.GHashLenResponse, error) {
	length, err := s.dbh.HLen(req.Key)
	if err != nil {
		return &ghash.GHashLenResponse{}, err
	}
	resp := &ghash.GHashLenResponse{
		Length: int64(length),
	}
	return resp, nil
}

func (s *hash) HUpdate(ctx context.Context, req *ghash.GHashUpdateRequest) (*ghash.GHashUpdateResponse, error) {

	ok, err := updateValue(s, req.Key, req.Field, req)
	if err != nil {
		return &ghash.GHashUpdateResponse{Ok: ok}, err
	}

	return &ghash.GHashUpdateResponse{Ok: ok}, nil
}

func updateValue(s *hash, key string, field interface{}, r *ghash.GHashUpdateRequest) (bool, error) {
	switch r.Value.(type) {
	case *ghash.GHashUpdateRequest_StringValue:
		ok, err := s.dbh.HUpdate(key, field, r.GetStringValue())
		return ok, err
	case *ghash.GHashUpdateRequest_Int32Value:
		ok, err := s.dbh.HUpdate(key, field, r.GetInt32Value())
		return ok, err
	case *ghash.GHashUpdateRequest_Int64Value:
		ok, err := s.dbh.HUpdate(key, field, r.GetInt64Value())
		return ok, err
	case *ghash.GHashUpdateRequest_Float32Value:
		ok, err := s.dbh.HUpdate(key, field, r.GetFloat32Value())
		return ok, err
	case *ghash.GHashUpdateRequest_Float64Value:
		ok, err := s.dbh.HUpdate(key, field, r.GetFloat64Value())
		return ok, err
	case *ghash.GHashUpdateRequest_BoolValue:
		ok, err := s.dbh.HUpdate(key, field, r.GetBoolValue())
		return ok, err
	case *ghash.GHashUpdateRequest_BytesValue:
		ok, err := s.dbh.HUpdate(key, field, r.GetBytesValue())
		return ok, err
	default:
		return false, fmt.Errorf("unknown value type")
	}
}

func (s *hash) HIncrBy(ctx context.Context, req *ghash.GHashIncrByRequest) (*ghash.GHashIncrByResponse, error) {
	value, err := s.dbh.HIncrBy(req.Key, req.Field, req.Value)
	if err != nil {
		return &ghash.GHashIncrByResponse{}, err
	}

	resp := &ghash.GHashIncrByResponse{
		Value: value,
	}
	return resp, nil
}

func (s *hash) HIncrByFloat(ctx context.Context, req *ghash.GHashIncrByFloatRequest) (*ghash.GHashIncrByFloatResponse, error) {
	value, err := s.dbh.HIncrByFloat(req.Key, req.Field, req.Value)
	if err != nil {
		return &ghash.GHashIncrByFloatResponse{}, err
	}

	resp := &ghash.GHashIncrByFloatResponse{
		Value: value,
	}
	return resp, nil
}

func (s *hash) HDecrBy(ctx context.Context, req *ghash.GHashDecrByRequest) (*ghash.GHashDecrByResponse, error) {
	value, err := s.dbh.HDecrBy(req.Key, req.Field, req.Value)
	if err != nil {
		return &ghash.GHashDecrByResponse{}, err
	}

	resp := &ghash.GHashDecrByResponse{
		Value: value,
	}
	return resp, nil
}

func (s *hash) HStrLen(ctx context.Context, req *ghash.GHashStrLenRequest) (*ghash.GHashStrLenResponse, error) {
	length, err := s.dbh.HStrLen(req.Key, req.Field)
	if err != nil {
		return &ghash.GHashStrLenResponse{}, err
	}
	resp := &ghash.GHashStrLenResponse{
		Length: int64(length),
	}
	return resp, nil
}

func (s *hash) HMove(ctx context.Context, req *ghash.GHashMoveRequest) (*ghash.GHashMoveResponse, error) {
	ok, err := s.dbh.HMove(req.Key, req.Dest, req.Field)
	if err != nil {
		return &ghash.GHashMoveResponse{Ok: ok}, err
	}
	return &ghash.GHashMoveResponse{Ok: ok}, nil
}

func (s *hash) HSetNX(ctx context.Context, req *ghash.GHashSetNXRequest) (*ghash.GHashSetNXResponse, error) {
	// Check if the key already exists
	ok, err := setNXValue(s, req.Key, req.Field, req.Ttl, req)
	if err != nil {
		return &ghash.GHashSetNXResponse{Ok: ok}, err
	}

	return &ghash.GHashSetNXResponse{Ok: ok}, nil
}
func setNXValue(s *hash, key string, field interface{}, ttl int64, r *ghash.GHashSetNXRequest) (bool, error) {
	switch r.Value.(type) {
	case *ghash.GHashSetNXRequest_StringValue:
		ok, err := s.dbh.HSetNX(key, field, r.GetStringValue(), ttl)
		return ok, err
	case *ghash.GHashSetNXRequest_Int32Value:
		ok, err := s.dbh.HSetNX(key, field, r.GetInt32Value(), ttl)
		return ok, err
	case *ghash.GHashSetNXRequest_Int64Value:
		ok, err := s.dbh.HSetNX(key, field, r.GetInt64Value(), ttl)
		return ok, err
	case *ghash.GHashSetNXRequest_Float32Value:
		ok, err := s.dbh.HSetNX(key, field, r.GetFloat32Value(), ttl)
		return ok, err
	case *ghash.GHashSetNXRequest_Float64Value:
		ok, err := s.dbh.HSetNX(key, field, r.GetFloat64Value(), ttl)
		return ok, err
	case *ghash.GHashSetNXRequest_BoolValue:
		ok, err := s.dbh.HSetNX(key, field, r.GetBoolValue(), ttl)
		return ok, err
	case *ghash.GHashSetNXRequest_BytesValue:
		ok, err := s.dbh.HSetNX(key, field, r.GetBytesValue(), ttl)
		return ok, err
	default:
		return false, fmt.Errorf("unknown value type")
	}
}

func (s *hash) HType(ctx context.Context, req *ghash.GHashTypeRequest) (*ghash.GHashTypeResponse, error) {
	hashType, err := s.dbh.HTypes(req.Key, req.Field)
	if err != nil {
		return &ghash.GHashTypeResponse{}, err
	}
	return &ghash.GHashTypeResponse{Type: hashType}, nil
}
