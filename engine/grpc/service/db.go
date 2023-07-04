package service

import (
	"context"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/ByteStorage/FlyDB/lib/proto"
)

// Service is a grpc service for db
type Service struct {
	proto.FlyDBServiceServer
	db *engine.DB
}

// Put is a grpc service for put
func (s *Service) Put(ctx context.Context, req *proto.PutRequest) (*proto.PutResponse, error) {
	err := s.db.Put([]byte(req.Key), []byte(req.Value))
	if err != nil {
		return &proto.PutResponse{}, err
	}
	return &proto.PutResponse{Ok: true}, nil
}

// Get is a grpc service for get
func (s *Service) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	value, err := s.db.Get([]byte(req.Key))
	if err != nil {
		return &proto.GetResponse{}, err
	}
	return &proto.GetResponse{Value: string(value)}, nil
}

// Del is a grpc service for del
func (s *Service) Del(ctx context.Context, req *proto.DelRequest) (*proto.DelResponse, error) {
	err := s.db.Delete([]byte(req.Key))
	if err != nil {
		return &proto.DelResponse{}, err
	}
	return &proto.DelResponse{Ok: true}, nil
}

// Keys is a grpc service for keys
func (s *Service) Keys(ctx context.Context, req *proto.KeysRequest) (*proto.KeysResponse, error) {
	list := s.db.GetListKeys()
	keys := make([]string, len(list))
	for i, bytes := range list {
		keys[i] = string(bytes)
	}
	return &proto.KeysResponse{Keys: keys}, nil
}
