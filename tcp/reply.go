package tcp

import (
	"context"
	"github.com/qishenonly/flydb/protocol/tcpIF"
	"github.com/qishenonly/flydb/sync/boolAm"
	"net"
	"sync"
)

var _ tcpIF.Handler = (*TcpReplyHandler)(nil)

type TcpReplyHandler struct {
	activeConn sync.Map
	isClosed   boolAm.Boolean
}

func (t *TcpReplyHandler) Handle(ct context.Context, conn net.Conn) {
	//TODO implement me
	panic("implement me")
}

func (t *TcpReplyHandler) Close() error {
	//TODO implement me
	panic("implement me")
}
