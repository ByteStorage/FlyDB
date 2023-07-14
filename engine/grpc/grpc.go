package grpc

type Grpc interface {
	// RegisterGrpcServer register a grpc server
	RegisterGrpcServer(addr string, class string) error
	// GetGrpcServer returns a grpc server
	GetGrpcServer(addr string, class string) (Service, error)
	// RegisterGrpcClient register a grpc client
	RegisterGrpcClient(addr string, class string) error
	// GetGrpcClient returns a grpc client
	GetGrpcClient(addr string, class string) (Client, error)
}

type Service interface {
}

type Client interface {
}

type grpc struct {
	services map[string]Service
	clients  map[string]Client
}

func NewGrpc() Grpc {
	return &grpc{
		services: make(map[string]Service),
		clients:  make(map[string]Client),
	}
}

func (g *grpc) RegisterGrpcServer(addr string, class string) error {
	//TODO implement me
	panic("implement me")
}

func (g *grpc) GetGrpcServer(addr string, class string) (Service, error) {
	//TODO implement me
	panic("implement me")
}

func (g *grpc) RegisterGrpcClient(addr string, class string) error {
	//TODO implement me
	panic("implement me")
}

func (g *grpc) GetGrpcClient(addr string, class string) (Client, error) {
	//TODO implement me
	panic("implement me")
}
