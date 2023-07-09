package service

import (
	"context"
	"fmt"
	"github.com/ByteStorage/FlyDB/lib/proto/dbs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os/signal"
	"syscall"
	"time"
)

// IsGrpcServerRunning returns whether the grpc server is running
func (s *Service) IsGrpcServerRunning() bool {
	conn, err := grpc.Dial(s.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false
	}
	err = conn.Close()
	return err == nil
}

// StartServer starts a grpc server
func (s *Service) StartServer() {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic("tcp listen error: " + err.Error())
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

// StopServer stops a grpc server
func (s *Service) StopServer() {
	s.sig <- syscall.SIGINT
}

// Put is a grpc s for put
func (s *Service) Put(ctx context.Context, req *dbs.SetRequest) (*dbs.SetResponse, error) {
	fmt.Println("receive put request: key: ", req.Key, " value: ", req.GetValue(), " duration: ", time.Duration(req.Expire))
	var err error
	switch req.Value.(type) {
	case *dbs.SetRequest_StringValue:
		err = s.db.Set(req.Key, req.GetStringValue(), time.Duration(req.Expire))
	case *dbs.SetRequest_Int32Value:
		err = s.db.Set(req.Key, req.GetInt32Value(), time.Duration(req.Expire))
	case *dbs.SetRequest_Int64Value:
		err = s.db.Set(req.Key, req.GetInt64Value(), time.Duration(req.Expire))
	case *dbs.SetRequest_Float32Value:
		err = s.db.Set(req.Key, req.GetFloat32Value(), time.Duration(req.Expire))
	case *dbs.SetRequest_Float64Value:
		err = s.db.Set(req.Key, req.GetFloat64Value(), time.Duration(req.Expire))
	case *dbs.SetRequest_BoolValue:
		err = s.db.Set(req.Key, req.GetBoolValue(), time.Duration(req.Expire))
	case *dbs.SetRequest_BytesValue:
		err = s.db.Set(req.Key, req.GetBytesValue(), time.Duration(req.Expire))
	default:
		err = fmt.Errorf("unknown value type")
	}
	if err != nil {
		return &dbs.SetResponse{}, err
	}
	return &dbs.SetResponse{Ok: true}, nil
}

// Get is a grpc s for get
func (s *Service) Get(ctx context.Context, req *dbs.GetRequest) (*dbs.GetResponse, error) {
	value, err := s.db.Get(req.Key)
	if err != nil {
		return &dbs.GetResponse{}, err
	}
	resp := &dbs.GetResponse{}
	switch v := value.(type) {
	case string:
		resp.Value = &dbs.GetResponse_StringValue{StringValue: v}
	case int32:
		resp.Value = &dbs.GetResponse_Int32Value{Int32Value: v}
	case int64:
		resp.Value = &dbs.GetResponse_Int64Value{Int64Value: v}
	case float32:
		resp.Value = &dbs.GetResponse_Float32Value{Float32Value: v}
	case float64:
		resp.Value = &dbs.GetResponse_Float64Value{Float64Value: v}
	case bool:
		resp.Value = &dbs.GetResponse_BoolValue{BoolValue: v}
	case []byte:
		resp.Value = &dbs.GetResponse_BytesValue{BytesValue: v}
	}
	return resp, nil
}

// Del is a grpc s for del
func (s *Service) Del(ctx context.Context, req *dbs.DelRequest) (*dbs.DelResponse, error) {
	err := s.db.Del(req.Key)
	if err != nil {
		return &dbs.DelResponse{}, err
	}
	return &dbs.DelResponse{Ok: true}, nil
}
