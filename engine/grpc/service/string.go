package service

import (
	"context"
	"fmt"
	"github.com/ByteStorage/FlyDB/lib/proto/gstring"
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
	gstring.RegisterGStringServiceServer(server, s)
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
	err = s.dbs.Stop()
	if err != nil {
		fmt.Println("flydb stop error: ", err)
		return
	}
	fmt.Println("flydb stop success on ", s.Addr)

}

// StopServer stops a grpc server
func (s *Service) StopServer() {
	s.sig <- syscall.SIGINT
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
	panic("implement me")
}

func (s *Service) Expire(ctx context.Context, req *gstring.ExpireRequest) (*gstring.ExpireResponse, error) {
	err := s.dbs.Expire(req.Key, time.Duration(req.Expire))
	if err != nil {
		return &gstring.ExpireResponse{}, err
	}
	return &gstring.ExpireResponse{Ok: true}, nil
}

func (s *Service) Persist(ctx context.Context, req *gstring.PersistRequest) (*gstring.PersistResponse, error) {
	panic("implement me")
}
