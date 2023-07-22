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
	result, err := setValue(s, req.Key, req.Field, req)
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

func setValue(s *hash, key string, field interface{}, r *ghash.GHashSetRequest) (bool, error) {
	switch r.Value.(type) {
	case *ghash.GHashSetRequest_StringValue:
		ok, err := s.dbh.HSet(key, field, r.GetStringValue())
		return ok, err
	case *ghash.GHashSetRequest_Int32Value:
		ok, err := s.dbh.HSet(key, field, r.GetInt32Value())
		return ok, err
	case *ghash.GHashSetRequest_Int64Value:
		ok, err := s.dbh.HSet(key, field, r.GetInt64Value())
		return ok, err
	case *ghash.GHashSetRequest_Float32Value:
		ok, err := s.dbh.HSet(key, field, r.GetFloat32Value())
		return ok, err
	case *ghash.GHashSetRequest_Float64Value:
		ok, err := s.dbh.HSet(key, field, r.GetFloat64Value())
		return ok, err
	case *ghash.GHashSetRequest_BoolValue:
		ok, err := s.dbh.HSet(key, field, r.GetBoolValue())
		return ok, err
	case *ghash.GHashSetRequest_BytesValue:
		ok, err := s.dbh.HSet(key, field, r.GetBytesValue())
		return ok, err
	default:
		return false, fmt.Errorf("unknown value type")
	}
}
