package grpc

import (
	"context"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/ByteStorage/FlyDB/lib/proto"
)

type Service struct {
	proto.FlyDBServiceServer
	db *engine.DB
}

func (s *Service) Put(ctx context.Context, req *proto.PutRequest) (*proto.PutResponse, error) {
	err := s.db.Put([]byte(req.Key), []byte(req.Value))
	if err != nil {
		return &proto.PutResponse{}, err
	}
	return &proto.PutResponse{Ok: true}, nil
}

func (s *Service) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	value, err := s.db.Get([]byte(req.Key))
	if err != nil {
		return &proto.GetResponse{}, err
	}
	return &proto.GetResponse{Value: string(value)}, nil
}

func (s *Service) Del(ctx context.Context, req *proto.DelRequest) (*proto.DelResponse, error) {
	err := s.db.Delete([]byte(req.Key))
	if err != nil {
		return &proto.DelResponse{}, err
	}
	return &proto.DelResponse{Ok: true}, nil
}

func (s *Service) Keys(ctx context.Context, req *proto.KeysRequest) (*proto.KeysResponse, error) {
	list := s.db.GetListKeys()
	keys := make([]string, len(list))
	for i, bytes := range list {
		keys[i] = string(bytes)
	}
	return &proto.KeysResponse{Keys: keys}, nil
}
