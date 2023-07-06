package service

import (
	"context"
	"fmt"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/ByteStorage/FlyDB/lib/proto/dbs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Service is a grpc Service for db
type Service struct {
	dbs.FlyDBServiceServer
	Addr string // db server address
	db   *engine.DB
	sig  chan os.Signal
}

// NewService returns a new grpc Service
func NewService(addr string, db *engine.DB) *Service {
	return &Service{
		Addr: addr,
		db:   db,
		sig:  make(chan os.Signal),
	}
}

// IsGrpcServerRunning returns whether the grpc server is running
func (s *Service) IsGrpcServerRunning() bool {
	conn, err := grpc.Dial(s.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false
	}
	err = conn.Close()
	if err != nil {
		return false
	}
	return true
}

// StartServer starts a grpc server
func (s *Service) StartServer() {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic("tcp listen error: " + err.Error())
		return
	}
	server := grpc.NewServer()
	dbs.RegisterFlyDBServiceServer(server, s)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	go func() {
		err := server.Serve(listener)
		if err != nil {
			panic("db server start error: " + err.Error())
		}
	}()
	//wait for server start
	for {
		conn, err := grpc.Dial(s.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		err = conn.Close()
		if err != nil {
			continue
		}
		break
	}
	fmt.Println("flydb start success on ", s.Addr)
	// graceful shutdown
	signal.Notify(s.sig, syscall.SIGINT, syscall.SIGKILL)

	<-s.sig
}

func (s *Service) StopServer() {
	s.sig <- syscall.SIGINT
}

// Put is a grpc s for put
func (s *Service) Put(ctx context.Context, req *dbs.PutRequest) (*dbs.PutResponse, error) {
	err := s.db.Put([]byte(req.Key), []byte(req.Value))
	if err != nil {
		return &dbs.PutResponse{}, err
	}
	return &dbs.PutResponse{Ok: true}, nil
}

// Get is a grpc s for get
func (s *Service) Get(ctx context.Context, req *dbs.GetRequest) (*dbs.GetResponse, error) {
	value, err := s.db.Get([]byte(req.Key))
	if err != nil {
		return &dbs.GetResponse{}, err
	}
	return &dbs.GetResponse{Value: value}, nil
}

// Del is a grpc s for del
func (s *Service) Del(ctx context.Context, req *dbs.DelRequest) (*dbs.DelResponse, error) {
	err := s.db.Delete([]byte(req.Key))
	if err != nil {
		return &dbs.DelResponse{}, err
	}
	return &dbs.DelResponse{Ok: true}, nil
}

// Keys is a grpc s for keys
func (s *Service) Keys(ctx context.Context, req *dbs.KeysRequest) (*dbs.KeysResponse, error) {
	list := s.db.GetListKeys()
	keys := make([][]byte, len(list))
	for i, bytes := range list {
		keys[i] = bytes
	}
	return &dbs.KeysResponse{Keys: keys}, nil
}
