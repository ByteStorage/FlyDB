package tcp

import (
	"context"
	"fmt"
	"github.com/qishenonly/flydb/protocol/tcpIF"
	"net"
	"sync"
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
// closeChan is a channel to close tcp server
// when closeChan receive a signal, tcp server will close
func ListenAndServe(listener net.Listener, handler tcpIF.Handler, closeChan <-chan struct{}) {
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	go func() {
		// wait for close signal
		<-closeChan
		fmt.Println("tcp server close but is shutting down...")
		_ = listener.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()
	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		fmt.Println("accept new connection: ", conn.RemoteAddr().String())
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	wg.Wait()
}
