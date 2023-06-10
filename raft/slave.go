package cluster

import (
	"context"
	"github.com/ByteStorage/flydb/lib/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func (s *Slave) StartGrpcServer() {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	go func() {
		err := server.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()
	// graceful shutdown
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	<-sig
}

func (s *Slave) RegisterToMaster() {

}

func (s *Slave) SendHeartbeat() {

}

func (s *Slave) ListenLeader() {

}

func (s *Slave) Get(ctx context.Context, in *proto.SlaveGetRequest) (*proto.SlaveGetResponse, error) {

	val, err := s.DB.Get([]byte(in.Key))
	if err != nil {
		return &proto.SlaveGetResponse{}, err
	}
	return &proto.SlaveGetResponse{Value: string(val)}, nil
}

func (s *Slave) Set(ctx context.Context, in *proto.SlaveSetRequest) (*proto.SlaveSetResponse, error) {
	err := s.DB.Put([]byte(in.Key), []byte(in.Value))
	if err != nil {
		return &proto.SlaveSetResponse{}, err
	}
	return &proto.SlaveSetResponse{Ok: true}, nil
}

func (s *Slave) Del(ctx context.Context, in *proto.SlaveDelRequest) (*proto.SlaveDelResponse, error) {
	err := s.DB.Delete([]byte(in.Key))
	if err != nil {
		return &proto.SlaveDelResponse{}, err
	}
	return &proto.SlaveDelResponse{Ok: true}, nil
}

func (s *Slave) Keys(ctx context.Context, in *proto.SlaveKeysRequest) (*proto.SlaveKeysResponse, error) {
	list := s.DB.GetListKeys()
	keys := make([]string, len(list))
	for i, bytes := range list {
		keys[i] = string(bytes)
	}
	return &proto.SlaveKeysResponse{Keys: keys}, nil

}

func (s *Slave) Exists(ctx context.Context, in *proto.SlaveExistsRequest) (*proto.SlaveExistsResponse, error) {
	_, err := s.Get(ctx, &proto.SlaveGetRequest{Key: in.Key})
	if err != nil {
		return &proto.SlaveExistsResponse{Exists: false}, err
	}
	return &proto.SlaveExistsResponse{Exists: true}, nil
}

func (s *Slave) Expire(ctx context.Context, in *proto.SlaveExpireRequest) (*proto.SlaveExpireResponse, error) {
	panic("implement me")

}

func (s *Slave) TTL(ctx context.Context, in *proto.SlaveTTLRequest) (*proto.SlaveTTLResponse, error) {
	panic("implement me")
}

func (s *Slave) Heartbeat(ctx context.Context, in *proto.SlaveHeartbeatRequest) (*proto.SlaveHeartbeatResponse, error) {
	panic("implement me")
}
