package cluster

import (
	"context"
	"github.com/ByteStorage/FlyDB/lib/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	for range time.Tick(200 * time.Millisecond) {
		// connect with the currently known "leader"
		conn, err := grpc.Dial(s.Leader, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		// The current known "leader" will tell the slave whether it is still the leader
		response, err := proto.NewMasterGrpcServiceClient(conn).
			RegisterSlave(context.Background(), &proto.MasterRegisterSlaveRequest{Addr: s.Addr})
		if err != nil && response.Ok {
			break
		}
		zap.L().Error("register slave failed", zap.Error(err))
		continue
	}
}

func (s *Slave) SendHeartbeat() {
	for range time.Tick(3 * time.Second) {
		// connect with the currently known "leader"
		conn, err := grpc.Dial(s.Leader, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		// The current known "leader" will tell the slave whether it is still the leader
		response, err := proto.NewMasterGrpcServiceClient(conn).
			ReceiveHeartbeat(context.Background(), &proto.MasterHeartbeatRequest{Addr: s.Addr})
		if err != nil && response.Ok {
			continue
		}
		zap.L().Error("heartbeat failed", zap.Error(err))
	}
}

func (s *Slave) ListenLeader() {
	for range time.Tick(200 * time.Millisecond) {
		// connect with the currently known "leader"
		conn, err := grpc.Dial(s.Leader, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		// The current known "leader" will tell the slave whether it is still the leader
		response, err := proto.NewMasterGrpcServiceClient(conn).
			CurrentLeader(context.Background(), &proto.MasterCurrentLeaderRequest{})
		if err != nil {
			continue
		}
		s.Leader = response.Leader
	}
}

func (s *Slave) UpdateSlaveMessage() {

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
