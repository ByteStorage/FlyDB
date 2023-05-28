package tcp

import (
	"context"
	"fmt"
	"github.com/qishenonly/flydb/protocol/tcpIF"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

// ListenAndServeBySignal start tcp server by signal
func ListenAndServeBySignal(cfg *Config, handler tcpIF.Handler) error {
	closeChan := make(chan struct{})
	// listen system-level signal
	signalChan := make(chan os.Signal)
	// syscall.SIGHUP: terminal closed
	// syscall.SIGINT: ctrl + c
	// syscall.SIGTERM: kill
	// syscall.SIGQUIT: ctrl + \
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// receive signal
	// when receive a signal, send a signal to closeChan
	// closeChan will close tcp server
	go func() {
		s := <-signalChan
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	fmt.Println("tcp server start listen on: ", cfg.Address)

	// start tcp server
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
