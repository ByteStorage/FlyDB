package base

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine/grpc/service"
	"github.com/ByteStorage/FlyDB/lib/proto/ghash"
	"github.com/ByteStorage/FlyDB/lib/proto/glist"
	"github.com/ByteStorage/FlyDB/lib/proto/gset"
	"github.com/ByteStorage/FlyDB/lib/proto/gstring"
	"github.com/ByteStorage/FlyDB/lib/proto/gzset"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Base interface {
	StartGrpcServer()
	StopGrpcServer()
	IsGrpcServerStarted() bool
	RegisterService(base service.Base)
}

type base struct {
	addr   string
	sig    chan os.Signal
	server *grpc.Server
	base   []service.Base
}

func NewService(options config.Options, addr string) (Base, error) {
	baseService := &base{
		addr:   addr,
		sig:    make(chan os.Signal),
		server: grpc.NewServer(),
		base:   make([]service.Base, 0),
	}
	// start string structure service
	stringService, err := service.NewStringService(options)
	if err != nil {
		return nil, err
	}
	baseService.RegisterService(stringService)
	gstring.RegisterGStringServiceServer(baseService.server, stringService)

	// start hash structure service
	hashService, err := service.NewHashService(options)
	if err != nil {
		return nil, err
	}
	baseService.RegisterService(hashService)
	ghash.RegisterGHashServiceServer(baseService.server, hashService)

	listService, err := service.NewListService(options)
	if err != nil {
		return nil, err
	}
	baseService.RegisterService(listService)
	glist.RegisterGListServiceServer(baseService.server, listService)

	zsetService, err := service.NewZSetService(options)
	if err != nil {
		return nil, err
	}
	baseService.RegisterService(zsetService)
	gzset.RegisterGZSetServiceServer(baseService.server, zsetService)
	return baseService, nil
}

func (s *base) RegisterService(base service.Base) {
	s.base = append(s.base, base)
}

// StartGrpcServer starts a grpc server
func (s *base) StartGrpcServer() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic("tcp listen error: " + err.Error())
	}
	server := s.server
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	go func() {
		err := server.Serve(listener)
		if err != nil {
			panic("db server start error: " + err.Error())
		}
	}()
	//wait for server start
	for {
		conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		err = conn.Close()
		if err != nil {
			continue
		}
		break
	}
	fmt.Println("flydb start success on ", s.addr)
	// graceful shutdown
	signal.Notify(s.sig, syscall.SIGINT, syscall.SIGKILL)

	<-s.sig
	for _, base := range s.base {
		err := base.CloseDb()
		if err != nil {
			fmt.Println("flydb stop error: ", err)
			return
		}
	}
	fmt.Println("flydb stop success on ", s.addr)
}

func (s *base) StopGrpcServer() {
	s.sig <- syscall.SIGINT
}

// IsGrpcServerStarted returns whether the grpc server is running
func (s *base) IsGrpcServerStarted() bool {
	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false
	}
	err = conn.Close()
	return err == nil
}
