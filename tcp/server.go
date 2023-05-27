package tcp

import (
	"context"
	"fmt"
	"github.com/qishenonly/flydb/protocol/tcpIF"
	"net"
)

type Config struct {
	Address string
}

// ListenAndServeBySignal start tcp server by signal
func ListenAndServeBySignal(cfg *Config, handler tcpIF.Handler) error {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	fmt.Println("tcp server start listen on: ", cfg.Address)

	closeChan := make(chan struct{})
	ListenAndServe(listener, handler, closeChan)

	return nil
}

// ListenAndServe start tcp server
func ListenAndServe(listener net.Listener, handler tcpIF.Handler, closeChan <-chan struct{}) {
	ctx := context.Background()
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		fmt.Println("accept new connection: ", conn.RemoteAddr().String())
		go func() {
			handler.Handle(ctx, conn)
		}()
	}
}
