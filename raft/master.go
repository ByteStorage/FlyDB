package cluster

import (
	"context"
	"github.com/ByteStorage/FlyDB/lib/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (m *Master) ListenSlave(slaves []Slave) {
	//TODO implement me
	panic("implement me")
}

func (m *Master) WaitForLeader() {
	timeTick := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-timeTick.C:
			ch := m.Raft.LeaderCh()
			select {
			case isLeader := <-ch:
				if isLeader {
					return
				}
			default:
				continue
			}
		}
	}
}

func (m *Master) NewRaft() {
	//判断是新启动的master还是重启的master
}

func (m *Master) StartGrpcServer() {
	listener, err := net.Listen("tcp", m.Addr)
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

func (m *Master) ListenRequest() {

}

func (m *Master) Get(ctx context.Context, in *proto.MasterGetRequest) (*proto.MasterGetResponse, error) {
	panic("implement me")
}

func (m *Master) Set(ctx context.Context, in *proto.MasterSetRequest) (*proto.MasterSetResponse, error) {
	panic("implement me")
}

func (m *Master) Del(ctx context.Context, in *proto.MasterDelRequest) (*proto.MasterDelResponse, error) {
	panic("implement me")
}

func (m *Master) Keys(ctx context.Context, in *proto.MasterKeysRequest) (*proto.MasterKeysResponse, error) {
	panic("implement me")
}

func (m *Master) Scan(ctx context.Context, in *proto.MasterScanRequest) (*proto.MasterScanResponse, error) {
	panic("implement me")
}

func (m *Master) Expire(ctx context.Context, in *proto.MasterExpireRequest) (*proto.MasterExpireResponse, error) {
	panic("implement me")
}

func (m *Master) TTL(ctx context.Context, in *proto.MasterTTLRequest) (*proto.MasterTTLResponse, error) {
	panic("implement me")
}

func (m *Master) Ping(ctx context.Context, in *proto.MasterPingRequest) (*proto.MasterPingResponse, error) {
	panic("implement me")
}

func (m *Master) Shutdown(ctx context.Context, in *proto.MasterShutdownRequest) (*proto.MasterShutdownResponse, error) {
	panic("implement me")
}

func (m *Master) RegisterSlave(ctx context.Context, in *proto.MasterRegisterSlaveRequest) (*proto.MasterRegisterSlaveResponse, error) {
	panic("implement me")
}

func (m *Master) CurrentLeader(ctx context.Context, in *proto.MasterCurrentLeaderRequest) (*proto.MasterCurrentLeaderResponse, error) {
	addr, _ := m.Raft.LeaderWithID()
	return &proto.MasterCurrentLeaderResponse{Leader: string(addr)}, nil
}
