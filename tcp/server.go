package tcp

import (
	"github.com/qishenonly/flydb/protocol/tcpIF"
	"net"
)

type Config struct {
	Address string
}

// ListenAndServeBySignal start tcp server by signal
func ListenAndServeBySignal(cfg *Config, handler tcpIF.Handler) error {
	return nil
}

// ListenAndServe start tcp server
func ListenAndServe(listener net.Listener, handler tcpIF.Handler, closeChan <-chan struct{}) {

}
