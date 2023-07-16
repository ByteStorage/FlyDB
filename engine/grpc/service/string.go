package service

import (
	"context"
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

// Put is a grpc s for put
func (s *Service) Put(ctx context.Context, req *gstring.SetRequest) (*gstring.SetResponse, error) {
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
func (s *Service) Get(ctx context.Context, req *gstring.GetRequest) (*gstring.GetResponse, error) {
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
func (s *Service) Del(ctx context.Context, req *gstring.DelRequest) (*gstring.DelResponse, error) {
	err := s.dbs.Del(req.Key)
	if err != nil {
		return &gstring.DelResponse{}, err
	}
	return &gstring.DelResponse{Ok: true}, nil
}

func (s *Service) Type(ctx context.Context, req *gstring.TypeRequest) (*gstring.TypeResponse, error) {
	panic("implement me")
}

func (s *Service) StrLen(ctx context.Context, req *gstring.StrLenRequest) (*gstring.StrLenResponse, error) {
	panic("implement me")
}

func (s *Service) GetSet(ctx context.Context, req *gstring.GetSetRequest) (*gstring.GetSetResponse, error) {
	panic("implement me")
}

func (s *Service) Append(ctx context.Context, req *gstring.AppendRequest) (*gstring.AppendResponse, error) {
	panic("implement me")
}

func (s *Service) Incr(ctx context.Context, req *gstring.IncrRequest) (*gstring.IncrResponse, error) {
	panic("implement me")
}

func (s *Service) IncrBy(ctx context.Context, req *gstring.IncrByRequest) (*gstring.IncrByResponse, error) {
	panic("implement me")
}

func (s *Service) IncrByFloat(ctx context.Context, req *gstring.IncrByFloatRequest) (*gstring.IncrByFloatResponse, error) {
	panic("implement me")
}

func (s *Service) Decr(ctx context.Context, req *gstring.DecrRequest) (*gstring.DecrResponse, error) {
	panic("implement me")
}

func (s *Service) DecrBy(ctx context.Context, req *gstring.DecrByRequest) (*gstring.DecrByResponse, error) {
	panic("implement me")
}

func (s *Service) Exists(ctx context.Context, req *gstring.ExistsRequest) (*gstring.ExistsResponse, error) {
	exists, err := s.dbs.Exists(req.Key)
	if err != nil {
		return nil, err
	}
	return &gstring.ExistsResponse{Exists: exists}, nil
}

func (s *Service) Expire(ctx context.Context, req *gstring.ExpireRequest) (*gstring.ExpireResponse, error) {
	err := s.dbs.Expire(req.Key, time.Duration(req.Expire)*time.Second)
	if err != nil {
		return &gstring.ExpireResponse{}, err
	}
	return &gstring.ExpireResponse{Ok: true}, nil
}

func (s *Service) Persist(ctx context.Context, req *gstring.PersistRequest) (*gstring.PersistResponse, error) {
	err := s.dbs.Persist(req.Key)
	if err != nil {
		return nil, err
	}
	return &gstring.PersistResponse{Ok: true}, nil
}

func (s *Service) MGet(ctx context.Context, req *gstring.MGetRequest) (*gstring.MGetResponse, error) {
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
